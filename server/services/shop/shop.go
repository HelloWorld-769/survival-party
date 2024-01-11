package shop

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"
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
		CreatedAt    time.Time
		ProductId    string
		Coins        int64
		Gems         int64
		Inventory    int64
		CurrencyType int64
		Price        int64
	}

	query := `SELECT uso.user_id,uso.purchased,uso.created_at, so.product_id,so.coins,so.gems,so.inventory,so.currency_type,so.price FROM user_special_offers uso
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
	}{
		ProductId:    specOfferRes.ProductId,
		Name:         "Starter Pack",
		Type:         4,
		CurrencyType: specOfferRes.CurrencyType,
		Price:        specOfferRes.Price,
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
