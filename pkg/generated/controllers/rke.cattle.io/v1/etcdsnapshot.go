/*
Copyright 2023 Rancher Labs, Inc.

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

	v1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ETCDSnapshotController interface for managing ETCDSnapshot resources.
type ETCDSnapshotController interface {
	generic.ControllerInterface[*v1.ETCDSnapshot, *v1.ETCDSnapshotList]
}

// ETCDSnapshotClient interface for managing ETCDSnapshot resources in Kubernetes.
type ETCDSnapshotClient interface {
	generic.ClientInterface[*v1.ETCDSnapshot, *v1.ETCDSnapshotList]
}

// ETCDSnapshotCache interface for retrieving ETCDSnapshot resources in memory.
type ETCDSnapshotCache interface {
	generic.CacheInterface[*v1.ETCDSnapshot]
}

type ETCDSnapshotStatusHandler func(obj *v1.ETCDSnapshot, status v1.ETCDSnapshotStatus) (v1.ETCDSnapshotStatus, error)

type ETCDSnapshotGeneratingHandler func(obj *v1.ETCDSnapshot, status v1.ETCDSnapshotStatus) ([]runtime.Object, v1.ETCDSnapshotStatus, error)

func RegisterETCDSnapshotStatusHandler(ctx context.Context, controller ETCDSnapshotController, condition condition.Cond, name string, handler ETCDSnapshotStatusHandler) {
	statusHandler := &eTCDSnapshotStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, generic.FromObjectHandlerToHandler(statusHandler.sync))
}

func RegisterETCDSnapshotGeneratingHandler(ctx context.Context, controller ETCDSnapshotController, apply apply.Apply,
	condition condition.Cond, name string, handler ETCDSnapshotGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &eTCDSnapshotGeneratingHandler{
		ETCDSnapshotGeneratingHandler: handler,
		apply:                         apply,
		name:                          name,
		gvk:                           controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterETCDSnapshotStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type eTCDSnapshotStatusHandler struct {
	client    ETCDSnapshotClient
	condition condition.Cond
	handler   ETCDSnapshotStatusHandler
}

func (a *eTCDSnapshotStatusHandler) sync(key string, obj *v1.ETCDSnapshot) (*v1.ETCDSnapshot, error) {
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

type eTCDSnapshotGeneratingHandler struct {
	ETCDSnapshotGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *eTCDSnapshotGeneratingHandler) Remove(key string, obj *v1.ETCDSnapshot) (*v1.ETCDSnapshot, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.ETCDSnapshot{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *eTCDSnapshotGeneratingHandler) Handle(obj *v1.ETCDSnapshot, status v1.ETCDSnapshotStatus) (v1.ETCDSnapshotStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ETCDSnapshotGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
