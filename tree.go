package gin

import (
	"errors"
	"fmt"
	"strings"
)

type node struct {
	path        string
	part        string
	isWild      bool
	middlewares []HandlerFunc
	handlers    map[string]HandlerFunc
	children    []*node
}

func (n *node) addMiddleware(path string, middlewares []HandlerFunc) {
	_ = n.insert("", path, parsePath(path), 0, middlewares, nil)
}

func (n *node) addRoute(method, path string, handler HandlerFunc) error {
	return n.insert(method, path, parsePath(path), 0, nil, handler)
}

func (n *node) getRoute(method, path string) (*[]HandlerFunc, HandlerFunc, map[string]string) {
	searchParts := parsePath(path)
	middlewares := make([]HandlerFunc, 0)
	searchNode := n.search(searchParts, 0, &middlewares)
	if searchNode != nil {
		params := make(map[string]string)
		parts := parsePath(searchNode.path)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' {
				key := "path"
				if len(part) > 1 {
					key = part[1:]
				}
				params[key] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return &middlewares, searchNode.handlers[method], params
	}
	return &middlewares, nil, nil
}

func (n *node) insert(method, path string, parts []string, height int, middlewares []HandlerFunc, handler HandlerFunc) error {
	if len(parts) == height {
		if handler != nil {
			if n.handlers == nil {
				n.handlers = make(map[string]HandlerFunc)
			}
			if _, ok := n.handlers[method]; ok {
				return errors.New(fmt.Sprintf("handler are already registered for path '%s'", n.path))
			}
			n.handlers[method] = handler
		} else {
			n.middlewares = append(n.middlewares, middlewares...)
		}
		n.path = path
		return nil
	}
	part := parts[height]
	err := n.checkWildPart(path, part)
	if err != nil {
		return err
	}
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: isWild(part)}
		n.children = append(n.children, child)
	}
	return child.insert(method, path, parts, height+1, middlewares, handler)
}

func (n *node) search(parts []string, height int, middlewares *[]HandlerFunc) *node {
	*middlewares = append(*middlewares, n.middlewares...)
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			value := *middlewares
			if length := len(value); length > 0 {
				*middlewares = append(*middlewares, value[:length-len(n.middlewares)]...)
			}
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1, middlewares)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	wildNodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part {
			nodes = append(nodes, child)
		} else if child.isWild {
			wildNodes = append(wildNodes, child)
		}
	}
	return append(nodes, wildNodes...)
}

func (n *node) checkWildPart(path, part string) error {
	paths := make([]string, 0)
	for _, child := range n.children {
		if child.part[0] == ':' && part[0] == ':' && child.part != part {
			child.getPath(&paths)
		}
		if (child.part[0] == '*' && part[0] != '*') || (child.part[0] != '*' && part[0] == '*') {
			child.getPath(&paths)
		}
	}
	if len(paths) > 0 {
		return errors.New(fmt.Sprintf("path '%s' conflicts with existing path '[%s]'", path, strings.Join(paths, ", ")))
	}
	return nil
}

func (n *node) getPath(paths *[]string) {
	if n != nil {
		if n.path != "" {
			*paths = append(*paths, n.path)
		}
		for _, child := range n.children {
			child.getPath(paths)
		}
	}
}

func isWild(part string) bool {
	return part[0] == ':' || part[0] == '*'
}

func parsePath(path string) []string {
	paths := strings.Split(path, "/")
	parts := make([]string, 0)
	for _, part := range paths {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}
