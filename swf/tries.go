package swf

import "strings"

// 两种特殊的路由层
// 1.参数匹配 如/d/:lang/go 可以匹配 /d/sherlock/go 与 /d/poirot/go
// 2.通配 如/d/*detective 可以匹配 /d/sherlock 与 /d/poirot/go

// 路由前缀树结点的结构
type Node struct {
	// 当前结点的路由值
	value string
	// 子结点的路由
	children []*Node
	// 是否为动态路由节点
	isDynamic bool
	// 结点的路由地址(也可用来判定是否为终止结点)
	pattern string
}

// 获取子结点中第一个匹配的结点
func (n *Node) GetFirstMatch(value string) *Node {
	// // 优先获取静态结点
	// for _, son := range n.children {
	// 	if son.value == value {
	// 		return son
	// 	}
	// }
	// // 动态结点的优先级在后
	// for _, son := range n.children {
	// 	if son.isDynamic {
	// 		return son
	// 	}
	// }
	// return nil
	// 不思考优先级的版本
	for _, son := range n.children {
		if son.value == value || son.isDynamic {
			return son
		}
	}
	return nil
}

// 获取子结点中所有匹配的结点
func (n *Node) GetAllMatch(value string) []*Node {
	res := make([]*Node, 0)
	for _, son := range n.children {
		if son.value == value || son.isDynamic {
			res = append(res, son)
		}
	}
	return res
}

// 添加新路由
func (n *Node) InsertTrie(pattern string, path []string, height int) {
	// 当前结点即为最终结点
	if height == len(path) {
		n.pattern = pattern
		return
	}
	// 取出当前层的路由结点并在当前结点的子结点中进行查找
	cur := path[height]
	next := n.GetFirstMatch(cur)
	if next != nil {
		// // 如果匹配到动态结点
		// if next.isDynamic {
		// 	curDynamic := cur[0] == '0' || cur[0] == '*'
		// 	// 如果当前路由层确实为动态层则向下匹配
		// 	if curDynamic {
		// 		next.InsertTrie(pattern, path, height+1)
		// 	} else {
		// 		// 如果当前不是动态层那么根据 静态 > 动态 > 通配 则需要新建立静态结点
		// 		son := &Node{
		// 			value:     cur,
		// 			children:  make([]*Node, 0),
		// 			isDynamic: false,
		// 		}
		// 		n.children = append(n.children, son)
		// 		son.InsertTrie(pattern, path, height+1)
		// 	}
		// } else {
		// 	// 是静态层则直接向下匹配
		// 	next.InsertTrie(pattern, path, height+1)
		// }
		// 不思考优先级的版本
		next.InsertTrie(pattern, path, height+1)
	} else {
		// 匹配结果为空则进行插入
		son := &Node{
			value:    cur,
			children: make([]*Node, 0),
			// 查看当前路由层是否为动态路由层
			isDynamic: cur[0] == ':' || cur[0] == '*',
		}
		n.children = append(n.children, son)
		// 向下继续插入
		son.InsertTrie(pattern, path, height+1)
	}
}

// 查询路由
func (n *Node) QueryTrie(path []string, height int) *Node {
	// 如果路由层数已经遍历完毕或者到达通配符*层则停止向下遍历
	if height == len(path) || strings.HasPrefix(n.value, "*") {
		// 是终点路由则返回(通配*一定是终点路由)
		if n.pattern != "" {
			return n
		} else {
			return nil
		}
	}
	cur := path[height]
	// 获取所有能匹配当前层的结点
	sons := n.GetAllMatch(cur)
	for _, v := range sons {
		// 对所有可能的结点都进行向下遍历
		result := v.QueryTrie(path, height+1)
		if result != nil {
			return result
		}
		// 目前查询还未满足 静态 > 动态 > 通配 需要调整...
	}
	return nil
}
