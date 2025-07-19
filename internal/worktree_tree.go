package internal

import (
	"fmt"
)

// WorktreeNode represents a node in the worktree tree
type WorktreeNode struct {
	Name     string
	Config   WorktreeConfig
	Parent   *WorktreeNode
	Children []*WorktreeNode
}

// WorktreeManager manages the worktree tree structure
type WorktreeManager struct {
	nodes map[string]*WorktreeNode
	roots []*WorktreeNode
}

// NewWorktreeManager creates a new WorktreeManager from parsed GBMConfig
func NewWorktreeManager(config *GBMConfig) (*WorktreeManager, error) {
	manager := &WorktreeManager{
		nodes: make(map[string]*WorktreeNode),
		roots: make([]*WorktreeNode, 0),
	}

	// First pass: create all nodes
	for name, wtConfig := range config.Worktrees {
		node := &WorktreeNode{
			Name:     name,
			Config:   wtConfig,
			Children: make([]*WorktreeNode, 0),
		}
		manager.nodes[name] = node
	}

	// Second pass: establish parent-child relationships
	for name, node := range manager.nodes {
		if node.Config.MergeInto != "" {
			// MergeInto contains the worktree name, not branch name
			parent, exists := manager.nodes[node.Config.MergeInto]
			if !exists {
				return nil, fmt.Errorf("worktree '%s' references non-existent merge_into target '%s'", name, node.Config.MergeInto)
			}
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		} else {
			// This is a root node
			manager.roots = append(manager.roots, node)
		}
	}

	// Check for circular dependencies
	if err := manager.detectCycles(); err != nil {
		return nil, err
	}

	return manager, nil
}

var ErrNoRootNodesFound = fmt.Errorf("no root nodes found (all nodes have merge_into)")
var ErrCircularDependency = fmt.Errorf("circular dependency detected")

// detectCycles uses DFS to detect circular dependencies in the worktree tree
func (wm *WorktreeManager) detectCycles() error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range wm.nodes {
		if !visited[name] {
			if wm.dfsHasCycle(name, visited, recStack) {
				return fmt.Errorf("%w: worktree '%s' is part of a circular dependency", ErrCircularDependency, name)
			}
		}
	}
	return nil
}

// dfsHasCycle performs depth-first search to detect cycles
func (wm *WorktreeManager) dfsHasCycle(nodeName string, visited, recStack map[string]bool) bool {
	visited[nodeName] = true
	recStack[nodeName] = true

	node := wm.nodes[nodeName]

	// Follow the merge_into relationship (parent relationship)
	if node.Config.MergeInto != "" {
		if !visited[node.Config.MergeInto] {
			if wm.dfsHasCycle(node.Config.MergeInto, visited, recStack) {
				return true
			}
		} else if recStack[node.Config.MergeInto] {
			// Found a back edge - cycle detected
			return true
		}
	}

	recStack[nodeName] = false
	return false
}

// GetNode returns a node by name
func (wm *WorktreeManager) GetNode(name string) *WorktreeNode {
	return wm.nodes[name]
}

// GetRoots returns all root nodes of the tree
func (wm *WorktreeManager) GetRoots() []*WorktreeNode {
	return wm.roots
}

// GetRoot returns the first root node (for backward compatibility)
// Use GetRoots() for multiple root scenarios
func (wm *WorktreeManager) GetRoot() *WorktreeNode {
	if len(wm.roots) > 0 {
		return wm.roots[0]
	}
	return nil
}

// GetAllNodes returns all nodes as a map
func (wm *WorktreeManager) GetAllNodes() map[string]*WorktreeNode {
	return wm.nodes
}

// GetAllDeepestLeafNodes returns the deepest leaf nodes from all root trees
func (wm *WorktreeManager) GetAllDeepestLeafNodes() []*WorktreeNode {
	var allDeepestLeaves []*WorktreeNode
	for _, root := range wm.roots {
		deepestLeaves := root.GetDeepestLeafNodes()
		allDeepestLeaves = append(allDeepestLeaves, deepestLeaves...)
	}
	return allDeepestLeaves
}

// Node traversal methods

// GetParent returns the parent node (upstream in merge chain)
func (wn *WorktreeNode) GetParent() *WorktreeNode {
	return wn.Parent
}

// GetChildren returns all child nodes (downstream in merge chain)
func (wn *WorktreeNode) GetChildren() []*WorktreeNode {
	return wn.Children
}

// GetPath returns the path from this node to the root
func (wn *WorktreeNode) GetPath() []*WorktreeNode {
	path := make([]*WorktreeNode, 0)
	current := wn
	for current != nil {
		path = append([]*WorktreeNode{current}, path...)
		current = current.Parent
	}
	return path
}

// GetDepth returns the depth of this node (root = 0)
func (wn *WorktreeNode) GetDepth() int {
	depth := 0
	current := wn.Parent
	for current != nil {
		depth++
		current = current.Parent
	}
	return depth
}

// IsRoot returns true if this node is the root node
func (wn *WorktreeNode) IsRoot() bool {
	return wn.Parent == nil
}

// IsLeaf returns true if this node has no children
func (wn *WorktreeNode) IsLeaf() bool {
	return len(wn.Children) == 0
}

// WalkUp traverses up the tree and calls the provided function for each node
func (wn *WorktreeNode) WalkUp(fn func(*WorktreeNode) bool) {
	current := wn
	for current != nil {
		if !fn(current) {
			break
		}
		current = current.Parent
	}
}

// WalkDown traverses down the tree (depth-first) and calls the provided function for each node
func (wn *WorktreeNode) WalkDown(fn func(*WorktreeNode) bool) {
	if !fn(wn) {
		return
	}
	for _, child := range wn.Children {
		child.WalkDown(fn)
	}
}

// GetLeafNodes returns all leaf nodes (nodes with no children) in the subtree
func (wn *WorktreeNode) GetLeafNodes() []*WorktreeNode {
	var leaves []*WorktreeNode
	wn.WalkDown(func(node *WorktreeNode) bool {
		if node.IsLeaf() {
			leaves = append(leaves, node)
		}
		return true
	})
	return leaves
}

// GetDeepestLeafNodes returns the leaf nodes at the maximum depth
func (wn *WorktreeNode) GetDeepestLeafNodes() []*WorktreeNode {
	leaves := wn.GetLeafNodes()
	if len(leaves) == 0 {
		return leaves
	}

	maxDepth := 0
	for _, leaf := range leaves {
		if depth := leaf.GetDepth(); depth > maxDepth {
			maxDepth = depth
		}
	}

	var deepestLeaves []*WorktreeNode
	for _, leaf := range leaves {
		if leaf.GetDepth() == maxDepth {
			deepestLeaves = append(deepestLeaves, leaf)
		}
	}
	return deepestLeaves
}

// PrintTree prints the tree structure starting from this node
func (wn *WorktreeNode) PrintTree(indent string) {
	fmt.Printf("%s%s (%s) -> %s\n", indent, wn.Name, wn.Config.Branch, wn.Config.Description)
	for _, child := range wn.Children {
		child.PrintTree(indent + "  ")
	}
}
