package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/amirhosssein0/simple-operator/api/v1alpha1"
)

type MiniAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.amir.local,resources=miniapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.amir.local,resources=miniapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.amir.local,resources=miniapps/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *MiniAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1) Fetch MiniApp
	var mini appsv1alpha1.MiniApp
	if err := r.Get(ctx, req.NamespacedName, &mini); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if mini.Spec.Image == "" {
		err := fmt.Errorf("MiniApp %s/%s: spec.image is required", mini.Namespace, mini.Name)
		logger.Error(err, "validation error")
		return ctrl.Result{}, err
	}

	var replicas int32 = 1
	if mini.Spec.Replicas != nil {
		replicas = *mini.Spec.Replicas
	}
	var port int32 = 8080
	if mini.Spec.Port != nil {
		port = *mini.Spec.Port
	}

	labels := map[string]string{"app": mini.Name}
	depName := mini.Name

	var dep appsv1.Deployment
	err := r.Get(ctx, types.NamespacedName{Name: depName, Namespace: mini.Namespace}, &dep)
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{}, err
	}

	desired := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      depName,
			Namespace: mini.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "app",
							Image: mini.Spec.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: port},
							},
						},
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(&mini, &desired, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if apierrors.IsNotFound(err) {
		logger.Info("Creating Deployment", "name", desired.Name)
		if err := r.Create(ctx, &desired); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	dep.Spec.Replicas = desired.Spec.Replicas
	dep.Spec.Selector = desired.Spec.Selector
	dep.Spec.Template = desired.Spec.Template

	logger.Info("Updating Deployment", "name", dep.Name)
	if err := r.Update(ctx, &dep); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MiniAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.MiniApp{}).
		Owns(&appsv1.Deployment{}).
		Named("miniapp").
		Complete(r)
}
