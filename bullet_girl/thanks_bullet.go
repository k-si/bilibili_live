package bullet_girl

import (
	"context"
	"github.com/k-si/bili_live/entity"
	"sync"
	"time"
)

// 检测到礼物，push [uname]->[giftName]->[cost]，number+1
// 每10s统计一次礼物，并进行感谢，礼物价值高于x元加一句大气

var thanksGiver *GiftThanksGiver

type GiftThanksGiver struct {
	giftTable map[string]map[string]int
	tableMu   sync.RWMutex
	giftChan  chan *entity.SendGiftText
}

func pushToGiftChan(g *entity.SendGiftText) {
	thanksGiver.giftChan <- g
}

func ThanksGift(ctx context.Context) {

	thanksGiver = &GiftThanksGiver{
		giftTable: make(map[string]map[string]int),
		tableMu:   sync.RWMutex{},
		giftChan:  make(chan *entity.SendGiftText, 1000),
	}

	var g *entity.SendGiftText
	var w = 10 * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			summarizeGift()
			t.Reset(w)
		case g = <-thanksGiver.giftChan:
			if thanksGiver.giftTable[g.Data.Uname] == nil {
				thanksGiver.giftTable[g.Data.Uname] = make(map[string]int)
			}
			thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName] += g.Data.Price
		}
	}
END:
}

func summarizeGift() {
	for name, m := range thanksGiver.giftTable {
		sumCost := 0
		for gift, cost := range m {

			// 名称长度适应
			zh := []rune(name)
			if len(zh) > 8 {
				PushToBulletSender("谢谢" + string(zh))
			} else {
				PushToBulletSender("谢谢" + string(zh) + "的" + gift + "～ 爱你")
			}

			// 计算打赏金额
			sumCost += cost

			// 感谢完后立刻清空map
			delete(m, gift)
		}

		// 总打赏高于x元，加一句大气
		if sumCost >= 200 { // 2元
			PushToBulletSender("老板大气大气")
		}
		delete(thanksGiver.giftTable, name)
	}
}
