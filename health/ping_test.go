package health

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ping", func() {
	var subject *Ping

	BeforeEach(func() {
		subject = NewPing(func() error {
			return nil
		}, time.Hour, 2, 3)
	})

	AfterEach(func() {
		subject.Stop()
	})

	It("should update health status", func() {
		Expect(subject.IsHealthy()).To(BeFalse())
		subject.update(true)
		Expect(subject.IsHealthy()).To(BeFalse())
		subject.update(true)
		Expect(subject.IsHealthy()).To(BeTrue())
		subject.update(true)
		Expect(subject.IsHealthy()).To(BeTrue())
		subject.update(false)
		Expect(subject.IsHealthy()).To(BeTrue())
		subject.update(false)
		Expect(subject.IsHealthy()).To(BeTrue())
		subject.update(false)
		Expect(subject.IsHealthy()).To(BeFalse())
		subject.update(false)
		Expect(subject.IsHealthy()).To(BeFalse())
		subject.update(true)
		Expect(subject.IsHealthy()).To(BeFalse())
		subject.update(true)
		Expect(subject.IsHealthy()).To(BeTrue())
	})

	It("should check periodically", func() {
		ping := NewPing(func() error {
			return nil
		}, time.Millisecond, 2, 3)
		defer ping.Stop()

		Expect(ping.IsHealthy()).To(BeFalse())
		Eventually(ping.IsHealthy, "20ms", "2ms").Should(BeTrue())
	})

})

func BenchmarkPing_IsHealthy(b *testing.B) {
	ping := NewPing(func() error { return nil }, time.Hour, 2, 3)
	defer ping.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ping.IsHealthy()
	}
}

func BenchmarkPing_update(b *testing.B) {
	ping := NewPing(func() error { return nil }, time.Hour, 2, 3)
	defer ping.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ping.update(true)
	}
}
