package shop

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type Res struct {
	SpecialOffer struct {
		ProductId string `json:"productId"`
		Name      string `json:"name"`
		Type      int64  `json:"type"`
		Data      []struct {
			RewardType int64 `json:"rewardType"`
			Quantity   int64 `json:"quantity"`
		} `json:"data"`
		CurrencyType int64 `json:"currencyType"`
		Price        int64 `json:"price"`
		ExpiresAt    int64 `json:"expiresAt"`
		IsAvailable  bool  `json:"isAvailable"`
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
	query = "SELECT * FROM shops "
	err = db.QueryExecutor(query, &shopDetails)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	fmt.Println("asjdhsad", utils.CalculateDays(specOfferRes.CreatedAt))
	var res Res

	res.SpecialOffer = struct {
		ProductId string "json:\"productId\""
		Name      string "json:\"name\""
		Type      int64  "json:\"type\""
		Data      []struct {
			RewardType int64 "json:\"rewardType\""
			Quantity   int64 "json:\"quantity\""
		} "json:\"data\""
		CurrencyType int64 "json:\"currencyType\""
		Price        int64 "json:\"price\""
		ExpiresAt    int64 `json:"expiresAt"`
		IsAvailable  bool  `json:"isAvailable"`
	}{
		ProductId:    specOfferRes.ProductId,
		Name:         "Starter Pack",
		Type:         4,
		CurrencyType: specOfferRes.CurrencyType,
		Price:        specOfferRes.Price,
		ExpiresAt:    int64(7 - utils.CalculateDays(specOfferRes.CreatedAt)),
		IsAvailable:  specOfferRes.IsAvailable,
		Data: []struct {
			RewardType int64 "json:\"rewardType\""
			Quantity   int64 "json:\"quantity\""
		}{
			{
				RewardType: utils.Coins,
				Quantity:   specOfferRes.Coins,
			}, {
				RewardType: utils.Gems,
				Quantity:   specOfferRes.Gems,
			}, {
				RewardType: utils.Inventory,
				Quantity:   specOfferRes.Inventory,
			}},
	}
	for _, data := range shopDetails {

		if data.RewardType == utils.Energy {
			res.Energy = Temp{
				Name: "energyShopView",
				Type: utils.Energy,
				Data: append(res.Energy.Data, data),
			}
		} else if data.RewardType == utils.Coins {
			res.Coins = Temp{
				Name: "coinShopView",
				Type: utils.Coins,
				Data: append(res.Coins.Data, data),
			}
		} else if data.RewardType == utils.Gems {
			res.Gems = Temp{
				Name: "gemShopView",
				Type: utils.Gems,
				Data: append(res.Gems.Data, data),
			}
		}

	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, res, ctx)

}

func BuyFromStoreService(ctx *gin.Context, userId string, input request.BuyStoreRequest) {
	var shopData model.Shop
	query := "SELECT * FROM shops WHERE product_id=?"
	err := db.QueryExecutor(query, &shopData, input.ProductId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	var userGameStats model.UserGameStats
	query = "SELECT * FROM user_game_stats WHERE user_id=?"
	err = db.QueryExecutor(query, &userGameStats, userId)
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
	} else {
		response.ShowResponse("Invalid currency type", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
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

	response.ShowResponse("Buy Success", utils.HTTP_OK, utils.SUCCESS, nil, ctx)

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
