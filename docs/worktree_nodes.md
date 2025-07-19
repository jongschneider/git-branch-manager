# Worktree YAML Parser
This Go implementation provides:
Core Structures:

WorktreeConfig: Represents the YAML configuration for each worktree
WorktreeNode: A tree node with parent/child relationships and traversal methods
WorktreeManager: Manages the entire tree structure

Key Features:

YAML Parsing: Uses gopkg.in/yaml.v3 to parse the configuration into Go structs
Tree Construction: Two-pass algorithm to build the tree:

First pass: Create all nodes
Second pass: Establish parent-child relationships


Traversal Methods:

GetParent() / GetChildren(): Direct navigation
GetPath(): Get full path from node to root
WalkUp() / WalkDown(): Visitor pattern traversal
GetDepth(): Calculate node depth


Utility Methods:

IsRoot() / IsLeaf(): Check node position
PrintTree(): Visualize the tree structure



Usage Example:
The code includes a complete example showing how to:

Parse the YAML configuration
Access specific nodes
Traverse up and down the tree
Get paths and depths

The last node in a tree is typically called a leaf node (a node with no children). In your merge chain context, these would be the "feature branches" or "development branches" that don't have anything merging into them.

GetLeafNodes() - Returns all leaf nodes (nodes with no children) in the subtree. In your merge chain context, these are the branches that nothing merges into - typically feature or development branches.
GetDeepestLeafNodes() - Returns only the leaf nodes that are at the maximum depth. These are the "furthest" branches from the root in your merge chain.

In your example configuration:

Root: master (depth 0)
Intermediate: preview (depth 1)
Leaf: production (depth 2) - this is both a leaf node and the deepest leaf node

The terminology:

Root node: The starting point (usually master or main)
Leaf nodes: End points with no children (feature/development branches)
Intermediate nodes: Nodes that have both parents and children

These methods are useful for operations like:

Finding all active development branches (GetLeafNodes())
Finding the most "downstream" branches that need attention (GetDeepestLeafNodes())
Determining merge order (start from deepest leaves and work toward root)


```go
package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

// WorktreeConfig represents the configuration for a single worktree
type WorktreeConfig struct {
	Branch      string `yaml:"branch"`
	Description string `yaml:"description"`
	MergeInto   string `yaml:"merge_into,omitempty"`
}

// Config represents the entire configuration file
type Config struct {
	Worktrees map[string]WorktreeConfig `yaml:"worktrees"`
}

// WorktreeNode represents a node in the worktree tree
type WorktreeNode struct {
	Name        string
	Config      WorktreeConfig
	Parent      *WorktreeNode
	Children    []*WorktreeNode
}

// WorktreeManager manages the worktree tree structure
type WorktreeManager struct {
	nodes map[string]*WorktreeNode
	root  *WorktreeNode
}

// NewWorktreeManager creates a new WorktreeManager from YAML configuration
func NewWorktreeManager(yamlData []byte) (*WorktreeManager, error) {
	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	manager := &WorktreeManager{
		nodes: make(map[string]*WorktreeNode),
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
			parent, exists := manager.nodes[node.Config.MergeInto]
			if !exists {
				return nil, fmt.Errorf("worktree '%s' references non-existent merge_into target '%s'", name, node.Config.MergeInto)
			}
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		} else {
			// This is a root node
			if manager.root != nil {
				return nil, fmt.Errorf("multiple root nodes found: '%s' and '%s'", manager.root.Name, name)
			}
			manager.root = node
		}
	}

	if manager.root == nil {
		return nil, fmt.Errorf("no root node found (node without merge_into)")
	}

	return manager, nil
}

// GetNode returns a node by name
func (wm *WorktreeManager) GetNode(name string) *WorktreeNode {
	return wm.nodes[name]
}

// GetRoot returns the root node of the tree
func (wm *WorktreeManager) GetRoot() *WorktreeNode {
	return wm.root
}

// GetAllNodes returns all nodes as a map
func (wm *WorktreeManager) GetAllNodes() map[string]*WorktreeNode {
	return wm.nodes
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

// Example usage
func main() {
	yamlConfig := `
worktrees:
  master:
    branch: master
    description: "Main production branch"
  preview:
    branch: production-2025-07-1
    description: "Blade Runner"
    merge_into: master
  production:
    branch: production-2025-05-1
    description: "Arrival"
    merge_into: preview
`

	manager, err := NewWorktreeManager([]byte(yamlConfig))
	if err != nil {
		log.Fatal(err)
	}

	// Print the entire tree
	fmt.Println("Worktree Tree Structure:")
	manager.GetRoot().PrintTree("")

	fmt.Println("\nTraversal Examples:")

	// Get a specific node and traverse up
	prodNode := manager.GetNode("production")
	if prodNode != nil {
		fmt.Printf("\nPath from '%s' to root:\n", prodNode.Name)
		path := prodNode.GetPath()
		for i, node := range path {
			fmt.Printf("%d. %s (%s)\n", i+1, node.Name, node.Config.Branch)
		}

		fmt.Printf("\nWalking up from '%s':\n", prodNode.Name)
		prodNode.WalkUp(func(node *WorktreeNode) bool {
			fmt.Printf("- %s: %s\n", node.Name, node.Config.Description)
			return true // continue walking
		})
	}

	// Walk down from root
	fmt.Println("\nWalking down from root:")
	manager.GetRoot().WalkDown(func(node *WorktreeNode) bool {
		indent := ""
		for i := 0; i < node.GetDepth(); i++ {
			indent += "  "
		}
		fmt.Printf("%s- %s (depth: %d)\n", indent, node.Name, node.GetDepth())
		return true
	})

	// Find leaf nodes
	fmt.Println("\nLeaf nodes (end of merge chains):")
	leafNodes := manager.GetRoot().GetLeafNodes()
	for _, leaf := range leafNodes {
		fmt.Printf("- %s (%s) at depth %d\n", leaf.Name, leaf.Config.Branch, leaf.GetDepth())
	}

	// Find deepest leaf nodes
	fmt.Println("\nDeepest leaf nodes:")
	deepestLeaves := manager.GetRoot().GetDeepestLeafNodes()
	for _, leaf := range deepestLeaves {
		fmt.Printf("- %s (%s) at depth %d\n", leaf.Name, leaf.Config.Branch, leaf.GetDepth())
	}
}
```
