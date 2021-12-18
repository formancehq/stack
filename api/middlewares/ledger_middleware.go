package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/numary/ledger/ledger"
	"log"
)

// LedgerMiddleware struct
type LedgerMiddleware struct {
	resolver *ledger.Resolver
}

// NewLedgerMiddleware
func NewLedgerMiddleware(
	resolver *ledger.Resolver,
) LedgerMiddleware {
	return LedgerMiddleware{
		resolver: resolver,
	}
}

// LedgerMiddleware
func (m *LedgerMiddleware) LedgerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("ledger")

		if name == "" {
			return
		}

		l, err := m.resolver.GetLedger(name)
		if err != nil {
			c.JSON(400, gin.H{
				"ok":  false,
				"err": err.Error(),
			})
		}
		defer func() {
			err := l.Close()
			if err != nil {
				log.Printf("error closing ledger: %s", err)
			}
		}()
		c.Set("ledger", l)

		c.Next()
	}
}
