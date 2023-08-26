// @author:戴林峰
// @date:2023/8/26
// @node:
package handler

//
//func Test_limit(t *testing.T) {
//	// 设置请求参数
//	targeter := vegeta.NewStaticTargeter(vegeta.Target{
//		Method: "POST",
//		URL:    "http://foobar.org/fnord/2",
//		Body:   []byte("Hello world!"),
//		Header: http.Header{
//			"Authorization": []string{"x67890"},
//			"Content-Type":  []string{"text/plain"},
//		},
//	})
//
//	// 设置压测参数
//	rate := vegeta.Rate{Freq: 300, Per: time.Second}
//	duration := 30 * time.Second
//	targets := make(chan *vegeta.Target, 300)
//	results := make(chan *vegeta.Result, 300)
//	var wg sync.WaitGroup
//	wg.Add(1)
//
//	// 并发执行请求
//	go func() {
//		defer wg.Done()
//		for tgt := range targets {
//			attacker := vegeta.NewAttacker()
//			res := attacker.Attack(tgt)
//			results <- res
//		}
//	}()
//
//	// 发送请求
//	go func() {
//		for res := range results {
//			if res.Error != "" {
//				fmt.Printf("Error: %v\n", res.Error)
//				continue
//			}
//			if res.Code != http.StatusOK {
//				fmt.Printf("Unexpected status code: %d\n", res.Code)
//				continue
//			}
//			body, err := ioutil.ReadAll(res.Body)
//			if err != nil {
//				fmt.Printf("Error reading response body: %v\n", err)
//				continue
//			}
//			fmt.Printf("Response body: %s\n", string(body))
//		}
//	}()
//
//	// 开始压测
//	for res := range vegeta.NewBatch(targeter, rate, duration) {
//		targets <- res
//	}
//
//	close(targets)
//	wg.Wait()
//	close(results)
//}
