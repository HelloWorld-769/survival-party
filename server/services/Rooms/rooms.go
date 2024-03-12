package rooms

import (
	"encoding/json"
	"fmt"
	"io"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

// TODO: Shift response messages to another file
func generateUUID() (generatedUUID uuid.UUID, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Handle the panic gracefully
			fmt.Println("Recovered from panic:", r)
			err = fmt.Errorf("panic during UUID generation: %v", r)
		}
	}()

	// Attempt to generate a new UUID
	generatedUUID = uuid.New()
	return generatedUUID, nil
}

type GameLeaveReq struct {
	ActorNr    int    `json:"ActorNr"`
	AppVersion string `json:"AppVersion"`
	AppId      string `json:"AppId"`
	GameId     string `json:"GameId"`
	IsInactive bool   `json:"IsInactive"`
	Reason     string `json:"Reason"`
	Region     string `json:"Region"`
	Type       string `json:"Type"`
	UserId     string `json:"UserId"`
	Nickname   string `json:"Nickname"`
}

type GameCloseReq struct {
	ActorCount int    `json:"ActorCount"`
	AppVersion string `json:"AppVersion"`
	AppId      string `json:"AppId"`
	GameId     string `json:"GameId"`
	Region     string `json:"Region"`
}

type GameJoinReq struct {
	ActorNr    int    `json:"ActorNr"`
	AppVersion string `json:"AppVersion"`
	AppId      string `json:"AppId"`
	GameId     string `json:"GameId"`
	Region     string `json:"Region"`
	Type       string `json:"Type"`
	UserId     string `json:"UserId"`
	Nickname   string `json:"Nickname"`
}

func GameClose(ctx *gin.Context) {

	bodybytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}
	var gameCloseReq GameCloseReq

	err = json.Unmarshal(bodybytes, &gameCloseReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}

	query := "UPDATE game_states SET is_running =false WHERE game_id=?"
	err = db.RawExecutor(query, gameCloseReq.GameId)
	if err != nil {
		fmt.Println("Error in updating the game state")
		return
	}

	//destroy this room from Database

	query = "delete from rooms where room_id =?"
	err = db.RawExecutor(query, gameCloseReq.GameId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}

	// query = "SELECT * FROM game_sates WHERE game_id=?"
	// err = db.QueryExecutor(query, nil, gameCloseReq.GameId)

	response.ShowResponse("room close succesfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

type GameStartReq struct {
	AppVersion string `json:"AppVersion"`
	AppId      string `json:"AppId"`
	// GameId     string    `json:"GameId"`
	Region    string    `json:"Region"`
	Type      string    `json:"Type"`
	RpcParams RPCParams `json:"RpcParams"`
	UserId    string    `json:"UserId"`
}

type RPCParams struct {
	RoomId string `json:"roomId"`
}

func GameStart(ctx *gin.Context) {
	bodybytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}
	var gameStartReq GameStartReq

	err = json.Unmarshal(bodybytes, &gameStartReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}
	fmt.Println("req body:--------->", string(bodybytes))
	fmt.Println("gameStartReq:--------> ", gameStartReq)

	//update the room state
	query := "update rooms set is_open=false where room_id=?"
	err = db.RawExecutor(query, gameStartReq.RpcParams.RoomId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("room close succesfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

func GameJoin(ctx *gin.Context) {

	bodybytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}
	var gameJoinReq GameJoinReq

	err = json.Unmarshal(bodybytes, &gameJoinReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}

	// var exists bool
	// query := "SELECT EXISTS (SELECT * FROm game_states WHERE game_id =? AND is_running=true)"
	// db.QueryExecutor(query, &exists, gameJoinReq.GameId)

	// if exists {
	// 	query := "UPDATE game_states set actor_nr=? where game_id =? AND is_running=true and user_id=?"
	// 	err = db.RawExecutor(query, gameJoinReq.ActorNr, gameJoinReq.GameId, gameJoinReq.UserId)
	// 	if err != nil {
	// 		fmt.Println("error is", err)
	// 		return
	// 	}
	// } else {

	if gameJoinReq.Nickname != "GameManager" {

		gameData := model.GameState{
			ActorNr:        gameJoinReq.ActorNr,
			GameId:         gameJoinReq.GameId,
			Time:           10,
			UserId:         gameJoinReq.UserId,
			TotalGames:     3,
			GamesCompleted: 0,
			IsDead:         false,
			IsZombie:       false,
			Xp:             1,
			Kills:          0,
			IsRunning:      true,
		}

		err = db.CreateRecord(&gameData)
		if err != nil {
			fmt.Println("Error in creating the record")
			return
		}
	}

	//join this user to the game room
	var userInRoom model.UsersInRooms
	userInRoom.RoomId = gameJoinReq.GameId
	userInRoom.UserId = gameJoinReq.UserId
	userInRoom.Actor_Nr = gameJoinReq.ActorNr

	err = db.CreateRecord(&userInRoom)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	//decrease the capacity of the room
	query := "update rooms set current_capacity=current_capacity-1 where room_id=?"
	err = db.RawExecutor(query, gameJoinReq.GameId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}

	response.ShowResponse("room join succesfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

func GameLeave(ctx *gin.Context) {

	fmt.Println("game leave called")

	bodybytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}
	var gameLeaveReq GameLeaveReq

	err = json.Unmarshal(bodybytes, &gameLeaveReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}
	fmt.Println("game leave req: ", gameLeaveReq)

	query := "delete from users_in_rooms where user_id=? and room_id=?"
	err = db.RawExecutor(query, gameLeaveReq.UserId, gameLeaveReq.GameId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	//fetch the room information
	//update the capacity of room in rooms table also
	query = "update rooms set current_capacity=current_capacity+1 where room_id=?"
	err = db.QueryExecutor(query, nil, gameLeaveReq.GameId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	//check whether it reaches it full capacity (if yes destroy the room)
	query = "delete from rooms where current_capacity=capacity"
	err = db.RawExecutor(query)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}

	response.ShowResponse("room leave succesfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)
}

// @Summary get room for connection
// @Description get room for new game
// @Tags Rooms
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Success 200 {object} response.Success "successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /get-room [get]
func GetRoom(ctx *gin.Context) {

	//create a post request to send to node server for actual room creation
	//create a uuid to send to node server of room

	newUUID, err := generateUUID()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	userId, exists := ctx.Get("userId")
	if !exists {
		response.ShowResponse("userId missing from ", utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	//check in the db whether a room is available with some capacity left

	var room model.Rooms
	query := "select * from rooms where current_capacity>=1 and is_open=true"

	err = db.QueryExecutor(query, &room)
	if err != nil {

		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}
	if room.RoomId != "" {

		//room already exists with some capacity
		//check wether this user is already connected to this room
		//add the user to the room and return the room id

		//TODO : Rewrite the query properly
		var exists bool
		query := "select exists(select *from users_in_rooms where user_id='" + userId.(string) + "' and room_id='" + room.RoomId + "')"
		err := db.QueryExecutor(query, &exists)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
		if exists {
			response.ShowResponse("user already in room", utils.HTTP_BAD_REQUEST, utils.FAILURE, room.RoomId, ctx)
			return
		}

		response.ShowResponse("room get successfully", utils.HTTP_OK, utils.SUCCESS, room.RoomId, ctx)
		return
	}
	apiURL := os.Getenv("RoomCreateURL") + userId.(string) + "/" + newUUID.String()
	fmt.Println("api url: ", apiURL)

	request, err := http.NewRequest("PUT", apiURL, nil)
	if err != nil {
		fmt.Println("errrr", err.Error())
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	fmt.Println("request--------->", request)

	// Set the Content-Type header
	request.Header.Set("Content-Type", "application/json")

	// Make the PUT request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)

		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return

	}
	fmt.Println("response: ", string(body))
	defer resp.Body.Close()

	if resp.StatusCode == 201 || resp.StatusCode == 200 {

		//create a new room in database with above uuid and userId

		//TODO: refactor it
		var room model.Rooms
		room.RoomId = newUUID.String()
		room.UserId = userId.(string)
		room.Capacity = 5
		room.Is_Open = true
		room.CurrentCapacity = 5

		err := db.CreateRecord(&room)
		if err != nil {

			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
		//update the room current capacity

		response.ShowResponse("room created successfully", utils.HTTP_OK, utils.SUCCESS, newUUID.String(), ctx)
		return
	} else {

		response.ShowResponse(string(body), int64(resp.StatusCode), utils.FAILURE, nil, ctx)
		return
	}

}

type GameDataHook struct {
	ActorNr int
	GameId  string
	UserID  string
}
type DefaultSuccessResponse struct {
	ResultCode int    `json:"ResultCode"`
	Message    string `json:"Message"`
}

func GameCreate(ctx *gin.Context) {

	//default success response
	var data GameDataHook
	body, _ := io.ReadAll(ctx.Request.Body)

	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error in unmarshalling the resposne from the hook")
		return
	}
	fmt.Println(string(body))

	gameData := model.GameState{
		ActorNr:        data.ActorNr,
		GameId:         data.GameId,
		Time:           10,
		UserId:         data.UserID,
		TotalGames:     3,
		GamesCompleted: 0,
		IsDead:         false,
		IsZombie:       false,
		Xp:             0,
		Kills:          0,
		IsRunning:      true,
	}

	err = db.CreateRecord(&gameData)
	if err != nil {
		fmt.Println("Error in creating the record")
		return
	}

	resp := &DefaultSuccessResponse{
		ResultCode: 0,
		Message:    "OK",
	}

	ctx.JSON(200, resp)

}
