package game

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var gameDifficulty string

//also implemented #3 and #4 bonus

func StartServer(playerAmount int) {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server started, waiting for %d players...\n", playerAmount)

	// Accept multiple amount of players

	conns := make([]net.Conn, playerAmount)
	var secretCode int

	//connects all the players
	for i := 0; i < playerAmount; i++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			i--
			continue
		}
		conns[i] = conn
		//welcomes and provides the entering player with the instructions
		writeToClient(conn, fmt.Sprintf("Welcome Player %d!\n", i+1))
		writeToClient(conn, `Game Instructions:
		- The goal is to guess a 4-digit secret number.
		- You have 20 seconds per turn.
		- Take turns guessing. First correct guess wins.
		- Type "exit" during your turn to leave the game.
		
		`)

		//first connected player chooses difficulty
		//if he's choice is unvalid he can retry. has a 20 second limit
		if i == 0 {

			deadline := time.Now().Add(20 * time.Second)

			for {
				writeToClient(conn, "Choose difficulty (easy, medium, hard):\n")
				buffer := make([]byte, 1024)
				conn.SetReadDeadline(deadline)

				n, err := conn.Read(buffer)
				if err != nil {
					writeToClient(conn, "⏱️ Timeout or error reading difficulty. Connection closed.\n")
					conn.Close()
					i--
					break
				}

				gameDifficulty = strings.ToLower(strings.TrimSpace(string(buffer[:n])))

				if gameDifficulty == "easy" || gameDifficulty == "medium" || gameDifficulty == "hard" {
					secretCode = GenerateSecretCode(gameDifficulty)
					log.Printf("Secret code generated with difficulty %s: %d", gameDifficulty, secretCode)
					break
				} else {
					writeToClient(conn, "Invalid input. Please try again.\n")
				}
			}
		}
	}

	startTurnBasedGame(conns, 20*time.Second)
}

func startTurnBasedGame(conns []net.Conn, turnTimeout time.Duration) {
	playerAmount := len(conns)

	for {
		secretCode := GenerateSecretCode(gameDifficulty)
		log.Printf("New round: secret code generated (difficulty: %s) = %d", gameDifficulty, secretCode)
		// could change per round if desired

		// Game loop (one round)
		var winner int = -1
		//-1 means no winner yet

		//this is a single game setup
		for {
			for i := 0; i < playerAmount; i++ {
				playerNum := i + 1
				currentConn := conns[i]

				writeToAll(conns, fmt.Sprintf("\n It's Player %d's turn!\n", playerNum))
				//checks if the game ended and if its an exit or a win
				ended, reason := handleTurn(currentConn, playerNum, secretCode, turnTimeout)

				if ended {
					winner = playerNum
					writeToAll(conns, fmt.Sprintf("Player %d %s. Game Over.\n", playerNum, reason))
					break
				}
			}
			if winner != -1 {
				break
			}
		}

		// Ask if they want to play again
		playAgain := true
		for i, conn := range conns {
			writeToClient(conn, "Do you want to play again? (yes/no):\n")

			buffer := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(20 * time.Second))
			n, err := conn.Read(buffer)

			//if timedout or didnt respond "yes" takes it as no
			if err != nil || strings.ToLower(strings.TrimSpace(string(buffer[:n]))) != "yes" {
				writeToClient(conn, "Thanks for playing. Goodbye!\n")
				playAgain = false
			} else {
				writeToClient(conn, "Great! Waiting for other players...\n")
			}
			_ = i
		}

		//if "no" is given or disconnect, ends it
		if !playAgain {
			for _, c := range conns {
				c.Close()
			}
			return
		}

		// if 'yes" was given, starts a new round
		writeToAll(conns, "\n New round starting...\n")
	}
}

func handleTurn(conn net.Conn, playerNum int, secretCode int, timeout time.Duration) (bool, string) {
	writeToClient(conn, fmt.Sprintf("Player %d, enter your 4-digit guess (you have %s):\n", playerNum, timeout))

	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(timeout))

	//if takes too long
	n, err := conn.Read(buffer)
	if err != nil {
		writeToClient(conn, "Timeout! Turn skipped.\n")
		return false, ""
	}

	guessStr := strings.TrimSpace(string(buffer[:n]))

	// Handle "exit" command
	if strings.ToLower(guessStr) == "exit" {
		writeToClient(conn, "You exited the game. Goodbye!\n")
		return true, "exited game" // signal to terminate game
	}

	guess, err := ValidateGuess(guessStr)
	if err != nil {
		writeToClient(conn, "Invalid guess: "+err.Error()+"\n")
		return false, ""
	}

	if guess == secretCode {
		timestamp := GenerateTimestampPrefix()
		writeToClient(conn, fmt.Sprintf("%s |  Correct guess!\n", timestamp))
		return true, "guessed the correct code"
	}

	writeToClient(conn, fmt.Sprintf("Player %d: Incorrect guess.\n", playerNum))

	return false, ""
}

func writeToClient(conn net.Conn, s string) {
	_, err := conn.Write([]byte(s))
	if err != nil {
		log.Printf("Error writing to client: %v", err)
		return
	}
}

func writeToAll(conns []net.Conn, msg string) {
	for _, c := range conns {
		writeToClient(c, msg)
	}
}
