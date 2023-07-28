package foo

import (
	"math/rand"
	"testing"
)

func TestF(t *testing.T) {
	if name := F(); name != "F" {
		t.Errorf("expected name is F but got %s", name)
	}
}

func TestQuickPow(t *testing.T) {
	t.Run("2^2 % 3", func(t *testing.T) {
		if res := QuickPow(2, 2, 3); res != 1 {
			t.Errorf("expected 1 but actual is %d", res)
		}
	})

	t.Run("100003^256 % 987654", func(t *testing.T) {
		if res := QuickPow(100003, 256, 987654); res != 301705 {
			t.Errorf("expected 301705 but actual is %d", res)
		}
	})
}

func BenchmarkQuickPow(b *testing.B) {
	//如果你需要做点预先配置的话可以使用b.ResetTimer()来重置计时器
	var mods int64 = 998244353
	var a []int64
	var s []int64
	for i := 0; i < b.N; i++ {
		a = append(a, rand.Int63())
		s = append(s, rand.Int63())
	}
	//从这一行以后才开始计算时间
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickPow(a[i], s[i], mods)
	}
}
