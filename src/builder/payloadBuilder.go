package builder

import (
	"fmt"
	"github.com/sitetester/tcp_service/src/common"
	"strconv"
	"strings"
)

type PayloadBuilder struct{}

func toInt(friends []string) []int {
	var friendsAsInts []int

	for _, i := range friends {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}

		friendsAsInts = append(friendsAsInts, j)
	}

	return friendsAsInts
}

func (b PayloadBuilder) BuildPayload() common.Payload {

	payload := common.Payload{
		UserId:  getUserId(),
		Friends: getFriends(),
	}

	return payload
}

func getUserId() int {

	var userId int

	fmt.Print("Enter userId: ")
	_, err := fmt.Scanf("%d", &userId)
	if err != nil {
		panic(err)
	}

	return userId
}

func getFriends() []int {

	var friendsStr string

	fmt.Print("Enter friends(comma separated integers): ")
	fmt.Scanf("%s", &friendsStr)

	if len(friendsStr) > 0 {
		friends := strings.Split(friendsStr, ",")
		return toInt(friends)
	}

	var friendsAsInts []int
	return friendsAsInts
}
