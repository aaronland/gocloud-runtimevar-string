package runtimevar

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronland/go-aws-session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	gc_runtimevar "gocloud.dev/runtimevar"
	_ "gocloud.dev/runtimevar/blobvar"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
	"net/url"
	"path/filepath"
	"strings"
)

func OpenString(ctx context.Context, url_str string) (string, error) {

	if url_str == "" {
		return "", errors.New("Invalid URL string")
	}
	
	parsed, err := url.Parse(url_str)

	if err != nil {
		return "", err
	}

	switch strings.ToUpper(parsed.Scheme) {
	case "AWSSECRETSMANAGER":

		query := parsed.Query()

		aws_region := query.Get("region")
		aws_creds := query.Get("credentials")

		if aws_region == "" {
			return "", errors.New("Missing parameter: region")
		}

		if aws_creds == "" {
			return "", errors.New("Missing parameter: credentials")
		}

		aws_session, err := session.NewSessionWithCredentials(aws_creds, aws_region)

		if err != nil {
			return "", err
		}

		manager := secretsmanager.New(aws_session)

		secret_id := filepath.Join(parsed.Host, parsed.Path)

		result, err := manager.GetSecretValue(&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secret_id),
		})

		if err != nil {
			return "", err
		}

		return *result.SecretString, nil

	default:

		path := fmt.Sprintf("%s?decoder=string", url_str)
		v, err := gc_runtimevar.OpenVariable(ctx, path)

		if err != nil {
			return "", err
		}

		defer v.Close()

		/*

			Latest is intended to be called per request, with the request context.
			It returns the latest good Snapshot of the variable value, blocking if
			no good value has ever been received. If ctx is Done, it returns the
			latest error indicating why no good value is available (not the ctx.Err()).
			You can pass an already-Done ctx to make Latest not block.

			https://godoc.org/gocloud.dev/runtimevar#Variable.Latest
		*/

		snapshot, err := v.Latest(ctx)

		if err != nil {
			return "", err
		}

		return snapshot.Value.(string), nil
	}

}
