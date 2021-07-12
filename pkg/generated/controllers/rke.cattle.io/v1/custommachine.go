/*
Copyright 2021 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type CustomMachineHandler func(string, *v1.CustomMachine) (*v1.CustomMachine, error)

type CustomMachineController interface {
	generic.ControllerMeta
	CustomMachineClient

	OnChange(ctx context.Context, name string, sync CustomMachineHandler)
	OnRemove(ctx context.Context, name string, sync CustomMachineHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() CustomMachineCache
}

type CustomMachineClient interface {
	Create(*v1.CustomMachine) (*v1.CustomMachine, error)
	Update(*v1.CustomMachine) (*v1.CustomMachine, error)
	UpdateStatus(*v1.CustomMachine) (*v1.CustomMachine, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.CustomMachine, error)
	List(namespace string, opts metav1.ListOptions) (*v1.CustomMachineList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CustomMachine, err error)
}

type CustomMachineCache interface {
	Get(namespace, name string) (*v1.CustomMachine, error)
	List(namespace string, selector labels.Selector) ([]*v1.CustomMachine, error)

	AddIndexer(indexName string, indexer CustomMachineIndexer)
	GetByIndex(indexName, key string) ([]*v1.CustomMachine, error)
}

type CustomMachineIndexer func(obj *v1.CustomMachine) ([]string, error)

type customMachineController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewCustomMachineController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) CustomMachineController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &customMachineController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromCustomMachineHandlerToHandler(sync CustomMachineHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.CustomMachine
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.CustomMachine))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *customMachineController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.CustomMachine))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateCustomMachineDeepCopyOnChange(client CustomMachineClient, obj *v1.CustomMachine, handler func(obj *v1.CustomMachine) (*v1.CustomMachine, error)) (*v1.CustomMachine, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *customMachineController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *customMachineController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *customMachineController) OnChange(ctx context.Context, name string, sync CustomMachineHandler) {
	c.AddGenericHandler(ctx, name, FromCustomMachineHandlerToHandler(sync))
}

func (c *customMachineController) OnRemove(ctx context.Context, name string, sync CustomMachineHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromCustomMachineHandlerToHandler(sync)))
}

func (c *customMachineController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *customMachineController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *customMachineController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *customMachineController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *customMachineController) Cache() CustomMachineCache {
	return &customMachineCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *customMachineController) Create(obj *v1.CustomMachine) (*v1.CustomMachine, error) {
	result := &v1.CustomMachine{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *customMachineController) Update(obj *v1.CustomMachine) (*v1.CustomMachine, error) {
	result := &v1.CustomMachine{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *customMachineController) UpdateStatus(obj *v1.CustomMachine) (*v1.CustomMachine, error) {
	result := &v1.CustomMachine{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *customMachineController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *customMachineController) Get(namespace, name string, options metav1.GetOptions) (*v1.CustomMachine, error) {
	result := &v1.CustomMachine{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *customMachineController) List(namespace string, opts metav1.ListOptions) (*v1.CustomMachineList, error) {
	result := &v1.CustomMachineList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *customMachineController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *customMachineController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.CustomMachine, error) {
	result := &v1.CustomMachine{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type customMachineCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *customMachineCache) Get(namespace, name string) (*v1.CustomMachine, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.CustomMachine), nil
}

func (c *customMachineCache) List(namespace string, selector labels.Selector) (ret []*v1.CustomMachine, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CustomMachine))
	})

	return ret, err
}

func (c *customMachineCache) AddIndexer(indexName string, indexer CustomMachineIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.CustomMachine))
		},
	}))
}

func (c *customMachineCache) GetByIndex(indexName, key string) (result []*v1.CustomMachine, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.CustomMachine, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.CustomMachine))
	}
	return result, nil
}

type CustomMachineStatusHandler func(obj *v1.CustomMachine, status v1.CustomMachineStatus) (v1.CustomMachineStatus, error)

type CustomMachineGeneratingHandler func(obj *v1.CustomMachine, status v1.CustomMachineStatus) ([]runtime.Object, v1.CustomMachineStatus, error)

func RegisterCustomMachineStatusHandler(ctx context.Context, controller CustomMachineController, condition condition.Cond, name string, handler CustomMachineStatusHandler) {
	statusHandler := &customMachineStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromCustomMachineHandlerToHandler(statusHandler.sync))
}

func RegisterCustomMachineGeneratingHandler(ctx context.Context, controller CustomMachineController, apply apply.Apply,
	condition condition.Cond, name string, handler CustomMachineGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &customMachineGeneratingHandler{
		CustomMachineGeneratingHandler: handler,
		apply:                          apply,
		name:                           name,
		gvk:                            controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterCustomMachineStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type customMachineStatusHandler struct {
	client    CustomMachineClient
	condition condition.Cond
	handler   CustomMachineStatusHandler
}

func (a *customMachineStatusHandler) sync(key string, obj *v1.CustomMachine) (*v1.CustomMachine, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type customMachineGeneratingHandler struct {
	CustomMachineGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *customMachineGeneratingHandler) Remove(key string, obj *v1.CustomMachine) (*v1.CustomMachine, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.CustomMachine{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *customMachineGeneratingHandler) Handle(obj *v1.CustomMachine, status v1.CustomMachineStatus) (v1.CustomMachineStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.CustomMachineGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}