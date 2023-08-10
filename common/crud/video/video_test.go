// @author:戴林峰
// @date:2023/8/5
// @node:
package video

import (
	"fmt"
	"testing"
)

func TestFindVideoListByUserId(t *testing.T) {
	id, _ := FindVideoListByUserId(1)
	fmt.Printf("%+v", id)
}
