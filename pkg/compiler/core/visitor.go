// Copyright 2016 Marapongo, Inc. All rights reserved.

package core

import (
	"github.com/marapongo/mu/pkg/ast"
	"github.com/marapongo/mu/pkg/diag"
)

// Visitor unifies all visitation patterns under a single interface.
type Visitor interface {
	Phase
	VisitWorkspace(w *ast.Workspace)
	VisitCluster(name string, cluster *ast.Cluster)
	VisitDependency(parent *ast.Workspace, ref ast.Ref, dep *ast.Dependency)
	VisitStack(stack *ast.Stack)
	VisitProperty(parent *ast.Stack, name string, prop *ast.Property)
	VisitServices(parent *ast.Stack, svcs *ast.Services)
	VisitService(pstack *ast.Stack, parent *ast.Services, name ast.Name, public bool, svc *ast.Service)
}

// NewInOrderVisitor wraps another Visitor and walks the tree in a deterministic order, deferring to another set of
// Visitor objects for pre- and/or post-actions.  Either pre or post may be nil.
func NewInOrderVisitor(pre Visitor, post Visitor) Visitor {
	return &inOrderVisitor{pre, post}
}

// inOrderVisitor simply implements the Visitor pattern as specified above.
//
// Note that we need to iterate all maps in a stable order (since Go's are unordered by default).  Sadly, this
// is rather verbose due to Go's lack of generics, reflectionless Keys() functions, and so on.
type inOrderVisitor struct {
	pre  Visitor
	post Visitor
}

var _ Visitor = &inOrderVisitor{} // compile-time assert that inOrderVisitor implements Visitor.

func (v *inOrderVisitor) Diag() diag.Sink {
	if v.pre != nil {
		return v.pre.Diag()
	}
	if v.post != nil {
		return v.post.Diag()
	}
	return nil
}

func (v *inOrderVisitor) VisitWorkspace(w *ast.Workspace) {
	if v.pre != nil {
		v.pre.VisitWorkspace(w)
	}

	for _, name := range ast.StableClusters(w.Clusters) {
		v.VisitCluster(name, w.Clusters[name])
	}
	for _, ref := range ast.StableDependencies(w.Dependencies) {
		v.VisitDependency(w, ref, w.Dependencies[ref])
	}

	if v.post != nil {
		v.post.VisitWorkspace(w)
	}
}

func (v *inOrderVisitor) VisitCluster(name string, cluster *ast.Cluster) {
	if v.pre != nil {
		v.pre.VisitCluster(name, cluster)
	}
	if v.post != nil {
		v.post.VisitCluster(name, cluster)
	}
}

func (v *inOrderVisitor) VisitDependency(parent *ast.Workspace, ref ast.Ref, dep *ast.Dependency) {
	if v.pre != nil {
		v.pre.VisitDependency(parent, ref, dep)
	}
	if v.post != nil {
		v.post.VisitDependency(parent, ref, dep)
	}
}

func (v *inOrderVisitor) VisitStack(stack *ast.Stack) {
	if v.pre != nil {
		v.pre.VisitStack(stack)
	}

	for _, name := range ast.StableProperties(stack.Properties) {
		v.VisitProperty(stack, name, stack.Properties[name])
	}
	v.VisitServices(stack, &stack.Services)

	if v.post != nil {
		v.post.VisitStack(stack)
	}
}

func (v *inOrderVisitor) VisitProperty(parent *ast.Stack, name string, prop *ast.Property) {
	if v.pre != nil {
		v.pre.VisitProperty(parent, name, prop)
	}
	if v.post != nil {
		v.post.VisitProperty(parent, name, prop)
	}
}

func (v *inOrderVisitor) VisitServices(parent *ast.Stack, svcs *ast.Services) {
	if v.pre != nil {
		v.pre.VisitServices(parent, svcs)
	}

	for _, name := range ast.StableServices(svcs.Private) {
		v.VisitService(parent, svcs, name, false, svcs.Private[name])
	}
	for _, name := range ast.StableServices(svcs.Public) {
		v.VisitService(parent, svcs, name, true, svcs.Public[name])
	}

	if v.post != nil {
		v.post.VisitServices(parent, svcs)
	}
}

func (v *inOrderVisitor) VisitService(pstack *ast.Stack, parent *ast.Services, name ast.Name,
	public bool, svc *ast.Service) {
	if v.pre != nil {
		v.pre.VisitService(pstack, parent, name, public, svc)
	}
	if v.post != nil {
		v.post.VisitService(pstack, parent, name, public, svc)
	}
}
