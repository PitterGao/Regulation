package amazonsChess

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"log"
	"time"
)

type Game struct {
	CurrentPlayer int                    `json:"current_player,omitempty"`
	CurrentState  *State                 `json:"current_state,omitempty"`
	Winner        int                    `json:"winner,omitempty"`
	Ai1Handler    func(*State) ChessMove `json:"ai_1_handler,omitempty"`
	Ai2Handler    func(*State) ChessMove `json:"ai_2_handler,omitempty"`
}

func NewGame(currentPlayer int) (*Game, error) {
	if currentPlayer != -1 && currentPlayer != 1 {
		return nil, errors.New("wrong currentPlayer(need -1 or 1)")
	}
	board := NewBoard()
	return &Game{
		CurrentPlayer: currentPlayer,
		CurrentState: &State{
			Board:         board,
			CurrentPlayer: currentPlayer,
		},
		Winner: 0,
	}, nil
}

// Reset a game, must call it before a round start if not by calling function start
func (g *Game) Reset(currentPlayer int) error {
	if currentPlayer != -1 && currentPlayer != 1 {
		return errors.New("wrong currentPlayer(need -1 or 1)")
	}
	g.CurrentPlayer = currentPlayer
	g.CurrentState = &State{
		Board:         NewBoard(),
		CurrentPlayer: currentPlayer,
	}
	g.Winner = 0
	return nil
}

func (g *Game) LogGenerate() ([]byte, error) {
	var oneLog Log
	if g.CurrentState.GameOver() != 0 {
		oneLog = Log{
			GameState: *g.CurrentState,
			Status:    1,
			Winner:    g.CurrentState.GameOver(),
		}
	} else {
		oneLog = Log{
			GameState: *g.CurrentState,
			Status:    0,
			Winner:    0,
		}
	}

	logJson, err := json.Marshal(oneLog)
	if err != nil {
		return nil, err
	}
	return logJson, nil
}

func (g *Game) GetMove(state *State) ChessMove {
	if g.CurrentPlayer == -1 {
		if g.Ai1Handler == nil {
			return ChessMove{}
		}
		return g.Ai1Handler(state)
	} else {
		if g.Ai2Handler == nil {
			return ChessMove{}
		}
		return g.Ai2Handler(state)
	}
}

func (g *Game) Start(isShow bool) [][]byte {
	var record [][]byte
	var logJson []byte
	var err error

	err = g.Reset(g.CurrentPlayer)
	if err != nil {
		log.Fatal(err)
	}

	logJson, err = g.LogGenerate()
	if err != nil {
		log.Fatal(err)
	}
	record = append(record, logJson)

	fmt.Print("\x1b7") // 保存光标位置 保存光标和Attrs <ESC> 7
	for g.CurrentState.GameOver() == 0 {
		var err error
		move := g.GetMove(g.CurrentState)
		if move.Equal(ChessMove{}) {
			g.CurrentState, _ = g.CurrentState.RandomMove()
		} else {
			g.CurrentState, err = g.CurrentState.StateMove(move)
			if err != nil {
				log.Fatal(err)
			}
		}
		g.CurrentPlayer = g.CurrentState.CurrentPlayer
		if isShow {

			fmt.Print("\x1b8")
			fmt.Print("\x1b[2k") // 清空当前行的内容 擦除线<ESC> [2K
			g.CurrentState.PrintState()

			time.Sleep(50 * time.Millisecond)
		}
		logJson, err = g.LogGenerate()
		if err != nil {
			log.Fatal(err)
		}
		record = append(record, logJson)
	}

	var playerStr string
	g.Winner = g.CurrentState.GameOver()
	if g.Winner == 1 {
		playerStr = color.New(color.FgHiRed).Sprintf("red")
	} else {
		playerStr = color.New(color.FgHiBlue).Sprintf("blue")
	}
	fmt.Printf("winner is: %s\n", playerStr)

	logJson, err = g.LogGenerate()
	if err != nil {
		log.Fatal(err)
	}
	record = append(record, logJson)

	return record
}

func NewBoard() []int {
	board := make([]int, 100)
	board[3] = -1
	board[6] = -1
	board[30] = -1
	board[39] = -1
	board[60] = 1
	board[69] = 1
	board[93] = 1
	board[96] = 1
	return board
}
