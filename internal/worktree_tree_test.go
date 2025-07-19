package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWorktreeManager(t *testing.T) {
	tests := []struct {
		name      string
		config    *GBMConfig
		expectErr func(t *testing.T, err error)
	}{
		{
			name: "Valid simple chain",
			config: &GBMConfig{
				Worktrees: map[string]WorktreeConfig{
					"master": {
						Branch:      "master",
						Description: "Main production branch",
					},
					"preview": {
						Branch:      "production-2025-07-1",
						Description: "Preview branch",
						MergeInto:   "master",
					},
				},
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Valid three-level chain",
			config: &GBMConfig{
				Worktrees: map[string]WorktreeConfig{
					"master": {
						Branch:      "master",
						Description: "Main production branch",
					},
					"preview": {
						Branch:      "production-2025-07-1",
						Description: "Preview branch",
						MergeInto:   "master",
					},
					"production": {
						Branch:      "production-2025-05-1",
						Description: "Production branch",
						MergeInto:   "preview",
					},
				},
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Invalid merge target",
			config: &GBMConfig{
				Worktrees: map[string]WorktreeConfig{
					"master": {
						Branch:      "master",
						Description: "Main production branch",
					},
					"preview": {
						Branch:      "production-2025-07-1",
						Description: "Preview branch",
						MergeInto:   "nonexistent-branch",
					},
				},
			},
			expectErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "references non-existent merge_into target")
			},
		},
		{
			name: "Multiple root nodes allowed",
			config: &GBMConfig{
				Worktrees: map[string]WorktreeConfig{
					"master": {
						Branch:      "master",
						Description: "Main production branch",
					},
					"develop": {
						Branch:      "develop",
						Description: "Development branch",
					},
				},
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "No root node - circular reference",
			config: &GBMConfig{
				Worktrees: map[string]WorktreeConfig{
					"preview": {
						Branch:      "production-2025-07-1",
						Description: "Preview branch",
						MergeInto:   "production",
					},
					"production": {
						Branch:      "production-2025-05-1",
						Description: "Production branch",
						MergeInto:   "preview",
					},
				},
			},
			expectErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "circular dependency detected")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewWorktreeManager(tt.config)
			tt.expectErr(t, err)
			if err == nil {
				assert.NotNil(t, manager)
				assert.NotNil(t, manager.GetRoot())
			} else {
				assert.Nil(t, manager)
			}
		})
	}
}

func TestWorktreeNode_BasicProperties(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
			"production": {
				Branch:      "production-2025-05-1",
				Description: "Production branch",
				MergeInto:   "preview",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	// Test root node
	root := manager.GetRoot()
	assert.Equal(t, "master", root.Name)
	assert.True(t, root.IsRoot())
	assert.False(t, root.IsLeaf())
	assert.Equal(t, 0, root.GetDepth())

	// Test intermediate node
	preview := manager.GetNode("preview")
	require.NotNil(t, preview)
	assert.Equal(t, "preview", preview.Name)
	assert.False(t, preview.IsRoot())
	assert.False(t, preview.IsLeaf())
	assert.Equal(t, 1, preview.GetDepth())

	// Test leaf node
	production := manager.GetNode("production")
	require.NotNil(t, production)
	assert.Equal(t, "production", production.Name)
	assert.False(t, production.IsRoot())
	assert.True(t, production.IsLeaf())
	assert.Equal(t, 2, production.GetDepth())
}

func TestWorktreeNode_Relationships(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
			"production": {
				Branch:      "production-2025-05-1",
				Description: "Production branch",
				MergeInto:   "preview",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	root := manager.GetRoot()
	preview := manager.GetNode("preview")
	production := manager.GetNode("production")

	// Test parent relationships
	assert.Nil(t, root.GetParent())
	assert.Equal(t, root, preview.GetParent())
	assert.Equal(t, preview, production.GetParent())

	// Test children relationships
	rootChildren := root.GetChildren()
	assert.Len(t, rootChildren, 1)
	assert.Equal(t, preview, rootChildren[0])

	previewChildren := preview.GetChildren()
	assert.Len(t, previewChildren, 1)
	assert.Equal(t, production, previewChildren[0])

	productionChildren := production.GetChildren()
	assert.Len(t, productionChildren, 0)
}

func TestWorktreeNode_GetPath(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
			"production": {
				Branch:      "production-2025-05-1",
				Description: "Production branch",
				MergeInto:   "preview",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	root := manager.GetRoot()
	preview := manager.GetNode("preview")
	production := manager.GetNode("production")

	// Test path from root
	rootPath := root.GetPath()
	assert.Len(t, rootPath, 1)
	assert.Equal(t, root, rootPath[0])

	// Test path from preview
	previewPath := preview.GetPath()
	assert.Len(t, previewPath, 2)
	assert.Equal(t, root, previewPath[0])
	assert.Equal(t, preview, previewPath[1])

	// Test path from production
	productionPath := production.GetPath()
	assert.Len(t, productionPath, 3)
	assert.Equal(t, root, productionPath[0])
	assert.Equal(t, preview, productionPath[1])
	assert.Equal(t, production, productionPath[2])
}

func TestWorktreeNode_WalkUp(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
			"production": {
				Branch:      "production-2025-05-1",
				Description: "Production branch",
				MergeInto:   "preview",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	production := manager.GetNode("production")

	var visited []string
	production.WalkUp(func(node *WorktreeNode) bool {
		visited = append(visited, node.Name)
		return true
	})

	expected := []string{"production", "preview", "master"}
	assert.Equal(t, expected, visited)
}

func TestWorktreeNode_WalkDown(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
			"production": {
				Branch:      "production-2025-05-1",
				Description: "Production branch",
				MergeInto:   "preview",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	root := manager.GetRoot()

	var visited []string
	root.WalkDown(func(node *WorktreeNode) bool {
		visited = append(visited, node.Name)
		return true
	})

	expected := []string{"master", "preview", "production"}
	assert.Equal(t, expected, visited)
}

func TestWorktreeNode_GetLeafNodes(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
			"production": {
				Branch:      "production-2025-05-1",
				Description: "Production branch",
				MergeInto:   "preview",
			},
			"feature": {
				Branch:      "feature-branch",
				Description: "Feature branch",
				MergeInto:   "master",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	root := manager.GetRoot()
	leaves := root.GetLeafNodes()

	assert.Len(t, leaves, 2)
	leafNames := make([]string, len(leaves))
	for i, leaf := range leaves {
		leafNames[i] = leaf.Name
	}
	assert.Contains(t, leafNames, "production")
	assert.Contains(t, leafNames, "feature")
}

func TestWorktreeNode_GetDeepestLeafNodes(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
			"production": {
				Branch:      "production-2025-05-1",
				Description: "Production branch",
				MergeInto:   "preview",
			},
			"feature": {
				Branch:      "feature-branch",
				Description: "Feature branch",
				MergeInto:   "master",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	root := manager.GetRoot()
	deepestLeaves := root.GetDeepestLeafNodes()

	// Only production should be at the deepest level (depth 2)
	assert.Len(t, deepestLeaves, 1)
	assert.Equal(t, "production", deepestLeaves[0].Name)
	assert.Equal(t, 2, deepestLeaves[0].GetDepth())
}

func TestWorktreeManager_GetNode(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	// Test existing nodes
	master := manager.GetNode("master")
	assert.NotNil(t, master)
	assert.Equal(t, "master", master.Name)

	preview := manager.GetNode("preview")
	assert.NotNil(t, preview)
	assert.Equal(t, "preview", preview.Name)

	// Test non-existent node
	nonexistent := manager.GetNode("nonexistent")
	assert.Nil(t, nonexistent)
}

func TestWorktreeManager_GetAllNodes(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main production branch",
			},
			"preview": {
				Branch:      "production-2025-07-1",
				Description: "Preview branch",
				MergeInto:   "master",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	allNodes := manager.GetAllNodes()
	assert.Len(t, allNodes, 2)
	assert.Contains(t, allNodes, "master")
	assert.Contains(t, allNodes, "preview")
}

func TestComplexTree(t *testing.T) {
	// Test with a more complex tree structure
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main branch",
			},
			"staging": {
				Branch:      "staging",
				Description: "Staging branch",
				MergeInto:   "master",
			},
			"preview": {
				Branch:      "preview",
				Description: "Preview branch",
				MergeInto:   "staging",
			},
			"dev": {
				Branch:      "develop",
				Description: "Development branch",
				MergeInto:   "staging",
			},
			"feature1": {
				Branch:      "feature-1",
				Description: "Feature 1",
				MergeInto:   "dev",
			},
			"feature2": {
				Branch:      "feature-2",
				Description: "Feature 2",
				MergeInto:   "dev",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	// Verify tree structure
	root := manager.GetRoot()
	assert.Equal(t, "master", root.Name)

	// Check depths
	assert.Equal(t, 0, manager.GetNode("master").GetDepth())
	assert.Equal(t, 1, manager.GetNode("staging").GetDepth())
	assert.Equal(t, 2, manager.GetNode("preview").GetDepth())
	assert.Equal(t, 2, manager.GetNode("dev").GetDepth())
	assert.Equal(t, 3, manager.GetNode("feature1").GetDepth())
	assert.Equal(t, 3, manager.GetNode("feature2").GetDepth())

	// Check leaf nodes
	leaves := root.GetLeafNodes()
	assert.Len(t, leaves, 3) // preview, feature1, feature2

	// Check deepest leaf nodes
	deepestLeaves := root.GetDeepestLeafNodes()
	assert.Len(t, deepestLeaves, 2) // feature1, feature2 (both at depth 3)
	leafNames := make([]string, len(deepestLeaves))
	for i, leaf := range deepestLeaves {
		leafNames[i] = leaf.Name
	}
	assert.Contains(t, leafNames, "feature1")
	assert.Contains(t, leafNames, "feature2")
}

func TestWorktreeManager_MultipleRoots(t *testing.T) {
	config := &GBMConfig{
		Worktrees: map[string]WorktreeConfig{
			"master": {
				Branch:      "master",
				Description: "Main branch",
			},
			"develop": {
				Branch:      "develop",
				Description: "Development branch",
			},
			"feature1": {
				Branch:      "feature-1",
				Description: "Feature 1",
				MergeInto:   "develop",
			},
			"hotfix": {
				Branch:      "hotfix-urgent",
				Description: "Urgent hotfix",
				MergeInto:   "master",
			},
		},
	}

	manager, err := NewWorktreeManager(config)
	require.NoError(t, err)

	// Should have two root nodes
	roots := manager.GetRoots()
	assert.Len(t, roots, 2)

	rootNames := make([]string, len(roots))
	for i, root := range roots {
		rootNames[i] = root.Name
	}
	assert.Contains(t, rootNames, "master")
	assert.Contains(t, rootNames, "develop")

	// GetRoot() should return the first root for backward compatibility
	firstRoot := manager.GetRoot()
	assert.NotNil(t, firstRoot)
	assert.Contains(t, []string{"master", "develop"}, firstRoot.Name)

	// Test GetAllDeepestLeafNodes across all trees
	allDeepestLeaves := manager.GetAllDeepestLeafNodes()
	assert.Len(t, allDeepestLeaves, 2) // feature1 and hotfix

	leafNames := make([]string, len(allDeepestLeaves))
	for i, leaf := range allDeepestLeaves {
		leafNames[i] = leaf.Name
	}
	assert.Contains(t, leafNames, "feature1")
	assert.Contains(t, leafNames, "hotfix")
}
