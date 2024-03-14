package shop

import (
	"fmt"
	"io/ioutil"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v2"
)

type Res struct {
	SpecialOffer struct {
		ProductId    string                    `json:"productId"`
		Name         string                    `json:"name"`
		Type         int64                     `json:"type"`
		Data         []response.RewardResponse `json:"data"`
		CurrencyType int64                     `json:"currencyType"`
		Price        int64                     `json:"price"`
		ExpiresAt    int64                     `json:"expiresAt"`
		IsAvailable  bool                      `json:"isAvailable"`
	}
	Energy Temp `json:"Energy"`
	Coins  Temp `json:"Coins"`
	Gems   Temp `json:"Gems"`
}

type Temp struct {
	Name string       `json:"name"`
	Type int64        `json:"type"`
	Data []model.Shop `json:"data"`
}

func GetStoreService(ctx *gin.Context, userId string) {
	fmt.Println("User Id: ", userId)

	var specOfferRes struct {
		UserId       string
		Purchased    bool
		IsAvailable  bool
		CreatedAt    time.Time
		ProductId    string
		Coins        int64
		Gems         int64
		Inventory    int64
		CurrencyType int64
		Price        int64
	}

	query := `SELECT uso.user_id,uso.purchased,so.is_available,uso.created_at, so.product_id,so.coins,so.gems,so.inventory,so.currency_type,so.price FROM user_special_offers uso
	JOIN special_offers so ON so.id=uso.special_offer_id
	WHERE uso.user_id=?`

	err := db.QueryExecutor(query, &specOfferRes, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	fmt.Println("")

	var shopDetails []model.Shop
	query = "SELECT * FROM shops where popup=false "
	err = db.QueryExecutor(query, &shopDetails)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var res Res

	res.SpecialOffer = struct {
		ProductId    string                    "json:\"productId\""
		Name         string                    "json:\"name\""
		Type         int64                     "json:\"type\""
		Data         []response.RewardResponse "json:\"data\""
		CurrencyType int64                     "json:\"currencyType\""
		Price        int64                     "json:\"price\""
		ExpiresAt    int64                     `json:"expiresAt"`
		IsAvailable  bool                      `json:"isAvailable"`
	}{
		ProductId:    specOfferRes.ProductId,
		Name:         "Starter Pack",
		Type:         4,
		CurrencyType: specOfferRes.CurrencyType,
		Price:        specOfferRes.Price,
		ExpiresAt:    int64(7 - utils.CalculateDays(specOfferRes.CreatedAt)),
		IsAvailable:  specOfferRes.IsAvailable,
		Data: []response.RewardResponse{
			{
				RewardType: utils.Coins,
				Quantity:   specOfferRes.Coins,
			}, {
				RewardType: utils.Gems,
				Quantity:   specOfferRes.Gems,
			},
			// }, {
			// 	RewardType: utils.Inventory,
			// 	Quantity:   specOfferRes.Inventory,
			// }},
		},
	}
	for _, data := range shopDetails {

		if data.RewardType == utils.Energy {
			res.Energy = Temp{
				Name: "energyShopView",
				Type: 1,
				Data: append(res.Energy.Data, data),
			}
		} else if data.RewardType == utils.Coins {
			res.Coins = Temp{
				Name: "coinShopView",
				Type: 1,
				Data: append(res.Coins.Data, data),
			}
		} else if data.RewardType == utils.Gems {
			res.Gems = Temp{
				Name: "gemShopView",
				Type: 1,
				Data: append(res.Gems.Data, data),
			}
		}

	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, res, ctx)

}

func BuyFromStoreService(ctx *gin.Context, userId string, input request.BuyStoreRequest) {

	if !db.RecordExist("shops", input.ProductId, "product_id") {
		response.ShowResponse(utils.RECORD_NOT_FOUND, utils.HTTP_NOT_FOUND, utils.FAILURE, nil, ctx)
		return
	}

	var shopData model.Shop
	if input.Popup {
		query := "SELECT * FROM shops WHERE product_id=? and popup=true"
		err := db.QueryExecutor(query, &shopData, input.ProductId)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
	} else {
		query := "SELECT * FROM shops WHERE product_id=? and popup=false"
		err := db.QueryExecutor(query, &shopData, input.ProductId)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
	}
	var userGameStats model.UserGameStats
	query := "SELECT * FROM user_game_stats WHERE user_id=?"
	err := db.QueryExecutor(query, &userGameStats, userId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	if shopData.CurrencyType == utils.C_GEMS {
		if userGameStats.CurrentGems < shopData.Price {
			response.ShowResponse("No enough gems", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}
		userGameStats.CurrentGems -= shopData.Price
	} else if shopData.CurrencyType == utils.C_COINS {
		if userGameStats.CurrentCoins < shopData.Price {
			response.ShowResponse("No enough coins", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}
		userGameStats.CurrentCoins -= shopData.Price
	} else if shopData.CurrencyType == utils.C_MONEY {

		//hit the google api

		jsonKeyFile := "server/survival.json"

		// Load the service account key file
		jsonKey, err := ioutil.ReadFile(jsonKeyFile)
		if err != nil {
			response.ShowResponse("Error reading file"+err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
		if err != nil {
			response.ShowResponse("Error"+err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		client := conf.Client(ctx)

		service, err := androidpublisher.New(client)
		if err != nil {
			response.ShowResponse("Error in making client"+err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		// Verify the purchase token
		resp, err := service.Purchases.Products.Get(utils.PACKAGE_NAME, shopData.ProductId, input.Token).Do()
		if err != nil {
			response.ShowResponse("Error in getting the data from api"+err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		fmt.Printf("Purchase state: %v\n", resp.PurchaseState)

	}

	if shopData.RewardType == utils.Energy {
		userGameStats.Energy += shopData.Quantity
	} else if shopData.RewardType == utils.Coins {
		userGameStats.CurrentCoins += shopData.Quantity
		userGameStats.TotalCoins += shopData.Quantity

	} else if shopData.RewardType == utils.Gems {
		userGameStats.CurrentGems += shopData.Quantity
		userGameStats.TotalGems += shopData.Quantity

	}

	//update player game stats
	updateFields := map[string]interface{}{
		"current_coins": userGameStats.CurrentCoins,
		"current_gems":  userGameStats.CurrentGems,
		"energy":        userGameStats.Energy,
		"total_coins":   userGameStats.TotalCoins,
		"total_gems":    userGameStats.TotalGems,
	}

	tx := db.BeginTransaction()

	err = tx.Model(&model.UserGameStats{}).Where("user_id=?", userId).Updates(updateFields).Error
	if err != nil {
		tx.Rollback()
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	err = tx.Commit().Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse("Buy Success", utils.HTTP_OK, utils.SUCCESS, struct {
		UpdatedCoins  int64 `json:"updatedCoins"`
		UpdatedGems   int64 `json:"updatedGems"`
		UpdatedEnergy int64 `json:"updatedEnergy"`
	}{
		UpdatedCoins:  userGameStats.CurrentCoins,
		UpdatedGems:   userGameStats.CurrentGems,
		UpdatedEnergy: userGameStats.Energy,
	}, ctx)

}

func GiveNewSpecialOffer() {
	var userId []string
	query := `SELECT user_id FROM user_special_offers
	WHERE created_at <= (now() - interval '7 days');`

	err := db.QueryExecutor(query, &userId)
	if err != nil {
		fmt.Println("Error ", err)
		return
	}

	query = `DELETE FROM user_special_offers
	WHERE created_at <= (now() - interval '7 days');`
	err = db.RawExecutor(query)
	if err != nil {
		fmt.Println("Error ", err)

		return
	}

	// Fetch offer IDs and shuffle them for better randomization
	var offerIds []string
	err = db.QueryExecutor("SELECT id FROM special_offers ORDER BY RANDOM()", &offerIds)
	if err != nil {
		fmt.Println("Error ", err)

		return
	}
	rand.Shuffle(len(offerIds), func(i, j int) { offerIds[i], offerIds[j] = offerIds[j], offerIds[i] })

	// Start a goroutine to handle offer assignment, recycling offers as needed
	offers := make(chan string, len(userId))
	go func() {
		i := 0
		for range userId {
			offers <- offerIds[i]
			i = (i + 1) % len(offerIds) // Recycle offers when reaching the end
		}
		close(offers)
	}()

	// Process user IDs concurrently
	for _, id := range userId {
		offerId := <-offers

		userStartPack := model.UserSpecialOffer{
			SpecialOfferId: offerId,
			UserId:         id,
			Purchased:      false,
		}

		err = db.CreateRecord(&userStartPack)
		if err != nil {
			fmt.Println("Error ", err)
			return
		}
	}

	fmt.Println("Sucessfully generated special rewards for that user ")
}

func GetPopupService(ctx *gin.Context, rewardId int64) {

	var result struct {
		Offers []struct {
			Id           string `json:"id"`
			CurrencyType int64  `json:"currencyType"`
			Price        int64  `json:"price"`
			Quantity     int64  `json:"quantity"`
		} `json:"offers"`
	}

	var storeDetails []model.Shop
	query := "SELECT * FROM shops WHERE reward_type=? AND popup=true order by currency_type DESC"
	err := db.QueryExecutor(query, &storeDetails, rewardId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	for _, data := range storeDetails {
		result.Offers = append(result.Offers, struct {
			Id           string `json:"id"`
			CurrencyType int64  `json:"currencyType"`
			Price        int64  `json:"price"`
			Quantity     int64  `json:"quantity"`
		}{
			Id:           data.ProductId,
			CurrencyType: data.CurrencyType,
			Price:        data.Price,
			Quantity:     data.Quantity,
		})
	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, result, ctx)

}
