package authorization

import (
	gocontext "context"
	"fmt"
	"strings"

	"github.com/kuadrant/authorino/pkg/auth"
	"github.com/kuadrant/authorino/pkg/context"
	"github.com/kuadrant/authorino/pkg/json"
	"github.com/kuadrant/authorino/pkg/log"

	kubeAuthz "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubeAuthzClient "k8s.io/client-go/kubernetes/typed/authorization/v1"
	"k8s.io/client-go/rest"
)

type kubernetesSubjectAccessReviewer interface {
	SubjectAccessReviews() kubeAuthzClient.SubjectAccessReviewInterface
}

func NewKubernetesAuthz(user json.JSONValue, groups []string, resourceAttributes *KubernetesAuthzResourceAttributes) (*KubernetesAuthz, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubernetesAuthz{
		User:               user,
		Groups:             groups,
		ResourceAttributes: resourceAttributes,
		authorizer:         k8sClient.AuthorizationV1(),
	}, nil
}

type KubernetesAuthzResourceAttributes struct {
	Namespace   json.JSONValue
	Group       json.JSONValue
	Resource    json.JSONValue
	Name        json.JSONValue
	SubResource json.JSONValue
	Verb        json.JSONValue
}

type KubernetesAuthz struct {
	User               json.JSONValue
	Groups             []string
	ResourceAttributes *KubernetesAuthzResourceAttributes

	authorizer kubernetesSubjectAccessReviewer
}

func (k *KubernetesAuthz) Call(pipeline auth.AuthPipeline, ctx gocontext.Context) (interface{}, error) {
	if err := context.CheckContext(ctx); err != nil {
		return false, err
	}

	authJSON := pipeline.GetAuthorizationJSON()
	jsonValueToStr := func(value json.JSONValue) string {
		return fmt.Sprintf("%s", value.ResolveFor(authJSON))
	}

	subjectAccessReview := kubeAuthz.SubjectAccessReview{
		Spec: kubeAuthz.SubjectAccessReviewSpec{
			User: jsonValueToStr(k.User),
		},
	}

	if k.ResourceAttributes != nil {
		resourceAttributes := k.ResourceAttributes

		subjectAccessReview.Spec.ResourceAttributes = &kubeAuthz.ResourceAttributes{
			Namespace:   jsonValueToStr(resourceAttributes.Namespace),
			Group:       jsonValueToStr(resourceAttributes.Group),
			Resource:    jsonValueToStr(resourceAttributes.Resource),
			Name:        jsonValueToStr(resourceAttributes.Name),
			Subresource: jsonValueToStr(resourceAttributes.SubResource),
			Verb:        jsonValueToStr(resourceAttributes.Verb),
		}
	} else {
		request := pipeline.GetHttp()

		subjectAccessReview.Spec.NonResourceAttributes = &kubeAuthz.NonResourceAttributes{
			Path: request.Path,
			Verb: strings.ToLower(request.Method),
		}
	}

	if len(k.Groups) > 0 {
		subjectAccessReview.Spec.Groups = k.Groups
	}

	log.FromContext(ctx).WithName("kubernetesauthz").V(1).Info("calling kubernetes subject access review api", "subjectaccessreview", subjectAccessReview)

	if result, err := k.authorizer.SubjectAccessReviews().Create(ctx, &subjectAccessReview, metav1.CreateOptions{}); err != nil {
		return false, err
	} else {
		return parseSubjectAccessReviewResult(result)
	}
}

func parseSubjectAccessReviewResult(subjectAccessReview *kubeAuthz.SubjectAccessReview) (bool, error) {
	status := subjectAccessReview.Status
	if status.Allowed {
		return true, nil
	} else {
		reason := status.Reason
		if reason == "" {
			reason = "unknown reason"
		}
		return false, fmt.Errorf("Not authorized: %s", reason)
	}
}
