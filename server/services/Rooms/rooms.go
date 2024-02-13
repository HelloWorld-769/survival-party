package rooms

import (
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
		var userInRoom model.UsersInRooms
		userInRoom.RoomId = room.RoomId
		userInRoom.UserId = userId.(string)
		err = db.CreateRecord(&userInRoom)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		//update the room current capacity
		room.CurrentCapacity = room.CurrentCapacity - 1

		err = db.UpdateRecord(&room, room.RoomId, "room_id").Error
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		response.ShowResponse("room get successfully", utils.HTTP_OK, utils.SUCCESS, newUUID.String(), ctx)

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
	fmt.Println("response: ", string(body))
	defer resp.Body.Close()

	if resp.StatusCode == 201 || resp.StatusCode == 200 {

		//create a new room in database with above uuid and userId

		var room model.Rooms
		room.RoomId = newUUID.String()
		room.UserId = userId.(string)
		room.Capacity = 5
		room.Is_Open = true
		room.CurrentCapacity = 4

		err := db.CreateRecord(&room)
		if err != nil {

			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
		//create entry in userRoom table also

		var userInRoom model.UsersInRooms

		userInRoom.RoomId = newUUID.String()
		userInRoom.UserId = userId.(string)
		err = db.CreateRecord(&userInRoom)
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
