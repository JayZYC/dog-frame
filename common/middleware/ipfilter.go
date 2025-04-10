package middleware

import (
	"github.com/dog-frame/common/logger"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"sync"
)

// IPFilterMode IP过滤模式
type IPFilterMode int

const (
	// IPFilterModeWhitelist 白名单模式 - 只允许列表中的IP访问
	IPFilterModeWhitelist IPFilterMode = iota
	// IPFilterModeBlacklist 黑名单模式 - 拒绝列表中的IP访问
	IPFilterModeBlacklist
)

// IPFilterConfig IP过滤配置
type IPFilterConfig struct {
	// 过滤模式
	Mode IPFilterMode
	// IP列表
	IPs []string
	// IP网段列表
	CIDRs []string
	// 排除的路径
	ExcludedPaths []string
	// 自定义IP获取函数
	IPExtractor func(*gin.Context) string
}

// DefaultIPFilterConfig 返回默认的IP过滤配置
func DefaultIPFilterConfig() IPFilterConfig {
	return IPFilterConfig{
		Mode:          IPFilterModeBlacklist,
		IPs:           []string{},
		CIDRs:         []string{},
		ExcludedPaths: []string{"/api/health", "/api/metrics"},
		IPExtractor: func(c *gin.Context) string {
			return c.ClientIP()
		},
	}
}

// ipFilter IP过滤器
type ipFilter struct {
	config     IPFilterConfig
	ipSet      map[string]struct{}
	cidrBlocks []*net.IPNet
	mu         sync.RWMutex
}

// newIPFilter 创建新的IP过滤器
func newIPFilter(config IPFilterConfig) *ipFilter {
	filter := &ipFilter{
		config: config,
		ipSet:  make(map[string]struct{}),
	}

	// 初始化IP集合
	for _, ip := range config.IPs {
		filter.ipSet[ip] = struct{}{}
	}

	// 初始化CIDR网段
	for _, cidr := range config.CIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err == nil {
			filter.cidrBlocks = append(filter.cidrBlocks, ipNet)
		}
	}

	return filter
}

// isAllowed 检查IP是否被允许
func (f *ipFilter) isAllowed(ip string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 检查IP是否在集合中
	_, inSet := f.ipSet[ip]

	// 检查IP是否在CIDR网段中
	inCIDR := false
	if !inSet && len(f.cidrBlocks) > 0 {
		parsedIP := net.ParseIP(ip)
		if parsedIP != nil {
			for _, cidr := range f.cidrBlocks {
				if cidr.Contains(parsedIP) {
					inCIDR = true
					break
				}
			}
		}
	}

	// 根据过滤模式决定是否允许
	switch f.config.Mode {
	case IPFilterModeWhitelist:
		return inSet || inCIDR
	case IPFilterModeBlacklist:
		return !inSet && !inCIDR
	default:
		return true
	}
}

// AddIP 添加IP到过滤器
func (f *ipFilter) AddIP(ip string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ipSet[ip] = struct{}{}
}

// RemoveIP 从过滤器中移除IP
func (f *ipFilter) RemoveIP(ip string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.ipSet, ip)
}

// AddCIDR 添加CIDR网段到过滤器
func (f *ipFilter) AddCIDR(cidr string) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.cidrBlocks = append(f.cidrBlocks, ipNet)
	return nil
}

// 全局IP过滤器
var globalIPFilter *ipFilter
var ipFilterOnce sync.Once

// getIPFilter 获取或创建IP过滤器
func getIPFilter(config IPFilterConfig) *ipFilter {
	ipFilterOnce.Do(func() {
		globalIPFilter = newIPFilter(config)
	})
	return globalIPFilter
}

// IPFilter 返回一个IP过滤中间件
func IPFilter() gin.HandlerFunc {
	return IPFilterWithConfig(DefaultIPFilterConfig())
}

// IPFilterWithConfig 返回一个使用自定义配置的IP过滤中间件
func IPFilterWithConfig(config IPFilterConfig) gin.HandlerFunc {
	filter := getIPFilter(config)

	return func(c *gin.Context) {
		// 检查是否是排除的路径
		path := c.Request.URL.Path
		for _, excludedPath := range config.ExcludedPaths {
			if strings.HasPrefix(path, excludedPath) {
				c.Next()
				return
			}
		}

		// 获取客户端IP
		clientIP := config.IPExtractor(c)

		// 检查IP是否被允许
		if !filter.isAllowed(clientIP) {
			logger.Warn(c, "IP access denied", "ip", clientIP, "path", path)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// GetIPFilter 获取全局IP过滤器实例
func GetIPFilter() *ipFilter {
	if globalIPFilter == nil {
		return getIPFilter(DefaultIPFilterConfig())
	}
	return globalIPFilter
}
