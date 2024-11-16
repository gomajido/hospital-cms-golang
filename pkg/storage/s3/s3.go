package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/constant"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
	"github.com/gomajido/hospital-cms-golang/pkg/storage"
)

// InitS3 initializes and returns an S3 instance configured with the provided S3 configuration.
func InitS3(s3Cfg *config.S3Config) *s3.Client {
	sessionToken := ""
	if s3Cfg.SessionToken != "" {
		sessionToken = s3Cfg.SessionToken
	}

	cfg, err := awsCfg.LoadDefaultConfig(context.TODO(),
		awsCfg.WithRegion(s3Cfg.Region),
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3Cfg.AccessKeyID, s3Cfg.SecretAccessKey, sessionToken)),
	)
	if err != nil {
		app_log.Fatalf("Error: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	if s3Cfg.EndpointUrl != "" {
		client = s3.NewFromConfig(cfg, func(options *s3.Options) {
			options.BaseEndpoint = aws.String(s3Cfg.EndpointUrl)
		})
	}
	return client
}

// AwsS3Usecase AwsS3Use case Storage implements the laravel Filesystem like interface using Amazon S3.
type AwsS3Usecase struct {
	Client *s3.Client
	Cfg    *config.S3Config
}

func NewAwsS3(s3Cfg *config.S3Config, client *s3.Client) storage.IStorageProviderRepository {
	return &AwsS3Usecase{
		Cfg:    s3Cfg,
		Client: client,
	}
}

// Exists checks if an object exists in the S3 bucket.
func (s *AwsS3Usecase) Exists(ctx context.Context, path string) bool {
	_, err := s.Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
	})
	return err == nil
}

// Get retrieves the content of an object from the S3 bucket as a string.
func (s *AwsS3Usecase) Get(ctx context.Context, path string) (string, error) {
	output, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return "", err
	}
	defer output.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	return buf.String(), err
}

// ReadStream retrieves the content of an object from the S3 bucket as a ReadCloser.
func (s *AwsS3Usecase) ReadStream(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

// Put stores content in an object within the S3 bucket.
func (s *AwsS3Usecase) Put(ctx context.Context, path string, contents interface{}, options ...interface{}) error {
	var body io.ReadSeeker

	switch v := contents.(type) {
	case string:
		body = bytes.NewReader([]byte(v))
	case []byte:
		body = bytes.NewReader(v)
	case io.Reader:
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(v); err != nil {
			return err
		}
		body = bytes.NewReader(buf.Bytes())
	default:
		return fmt.Errorf("unsupported content type")
	}

	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
		Body:   body,
	})
	return err
}

// WriteStream writes the content of a Reader to an object within the S3 bucket.
func (s *AwsS3Usecase) WriteStream(ctx context.Context, path string, reader io.Reader, options ...interface{}) error {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return err
	}
	return s.Put(ctx, path, buf.Bytes())
}

// GetVisibility retrieves the visibility setting of an object in the S3 bucket.
func (s *AwsS3Usecase) GetVisibility(ctx context.Context, path string) (string, error) {
	output, err := s.Client.GetObjectAcl(context.TODO(), &s3.GetObjectAclInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return "", err
	}

	for _, grant := range output.Grants {
		if grant.Grantee != nil && grant.Grantee.Type == types.TypeGroup && *grant.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
			if grant.Permission == types.PermissionRead {
				return constant.PUBLIC, nil
			}
		}
	}
	return constant.PRIVATE, nil
}

// SetVisibility sets the visibility of an object in the S3 bucket.
func (s *AwsS3Usecase) SetVisibility(ctx context.Context, path, visibility string) error {
	var acl string
	switch visibility {
	case "public":
		acl = "public-read"
	case "private":
		acl = "private"
	default:
		return fmt.Errorf("unsupported visibility type")
	}

	_, err := s.Client.PutObjectAcl(context.TODO(), &s3.PutObjectAclInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
		ACL:    types.ObjectCannedACL(acl),
	})
	return err
}

// Prepend prepends data to the beginning of an object's content in the S3 bucket.
func (s *AwsS3Usecase) Prepend(ctx context.Context, path, data string) error {
	contents, err := s.Get(ctx, path)
	if err != nil {
		return err
	}
	return s.Put(ctx, path, data+contents)
}

// Append appends data to the end of an object's content in the S3 bucket.
func (s *AwsS3Usecase) Append(ctx context.Context, path, data string) error {
	contents, err := s.Get(ctx, path)
	if err != nil {
		return err
	}
	return s.Put(ctx, path, contents+data)
}

// Delete removes one or more objects from the S3 bucket.
func (s *AwsS3Usecase) Delete(ctx context.Context, paths ...string) error {
	for _, path := range paths {
		_, err := s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(s.Cfg.Bucket),
			Key:    aws.String(path),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Copy copies an object from one location to another within the S3 bucket.
func (s *AwsS3Usecase) Copy(ctx context.Context, from, to string) error {
	_, err := s.Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(s.Cfg.Bucket),
		CopySource: aws.String(filepath.Join(s.Cfg.Bucket, from)),
		Key:        aws.String(to),
	})
	return err
}

// Move moves an object from one location to another within the S3 bucket.
func (s *AwsS3Usecase) Move(ctx context.Context, from, to string) error {
	if err := s.Copy(ctx, from, to); err != nil {
		return err
	}
	return s.Delete(ctx, from)
}

// Size retrieves the size of an object in the S3 bucket.
func (s *AwsS3Usecase) Size(ctx context.Context, path string) (int64, error) {
	output, err := s.Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return 0, err
	}
	return *output.ContentLength, nil
}

// LastModified retrieves the last modified timestamp of an object in the S3 bucket.
func (s *AwsS3Usecase) LastModified(ctx context.Context, path string) (int64, error) {
	output, err := s.Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return 0, err
	}
	return output.LastModified.Unix(), nil
}

// Files lists the files under a directory in the S3 bucket, optionally recursively.
func (s *AwsS3Usecase) Files(ctx context.Context, directory string, recursive bool) ([]string, error) {
	return s.listObjects(ctx, directory, false, recursive)
}

// AllFiles lists all files under a directory in the S3 bucket, including those in subdirectories.
func (s *AwsS3Usecase) AllFiles(ctx context.Context, directory string) ([]string, error) {
	return s.listObjects(ctx, directory, false, true)
}

// Directories lists the directories under a directory in the S3 bucket, optionally recursively.
func (s *AwsS3Usecase) Directories(ctx context.Context, directory string, recursive bool) ([]string, error) {
	return s.listObjects(ctx, directory, true, recursive)
}

// AllDirectories lists all directories under a directory in the S3 bucket, including those in subdirectories.
func (s *AwsS3Usecase) AllDirectories(ctx context.Context, directory string) ([]string, error) {
	return s.listObjects(ctx, directory, true, true)
}

// MakeDirectory creates a directory in the S3 bucket by storing an empty object with a trailing slash.
func (s *AwsS3Usecase) MakeDirectory(ctx context.Context, path string) error {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Cfg.Bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader([]byte("")),
	})
	return err
}

// DeleteDirectory deletes a directory in the S3 bucket by removing all its contained objects.
func (s *AwsS3Usecase) DeleteDirectory(ctx context.Context, directory string) error {
	files, err := s.AllFiles(ctx, directory)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := s.Delete(ctx, file); err != nil {
			return err
		}
	}
	return nil
}

// listObjects lists objects under a directory in the S3 bucket, optionally filtering for directories only and optionally recursively.
func (s *AwsS3Usecase) listObjects(_ context.Context, directory string, directoriesOnly bool, recursive bool) ([]string, error) {
	var result []string
	delimiter := "/"
	if recursive {
		delimiter = ""
	}

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.Cfg.Bucket),
		Prefix:    aws.String(directory),
		Delimiter: aws.String(delimiter),
	}

	for {
		output, err := s.Client.ListObjectsV2(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		for _, prefix := range output.CommonPrefixes {
			if directoriesOnly {
				result = append(result, *prefix.Prefix)
			}
		}

		for _, object := range output.Contents {
			if !directoriesOnly {
				result = append(result, *object.Key)
			}
		}

		if output.IsTruncated != nil && *output.IsTruncated {
			input.ContinuationToken = output.NextContinuationToken
		} else {
			break
		}
	}

	return result, nil
}

func (s *AwsS3Usecase) GenerateURL(ctx context.Context, path string, expires time.Duration) (*string, error) {
	presignClient := s3.NewPresignClient(s.Client)
	presignURL, err := presignClient.PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Cfg.Bucket),
			Key:    aws.String(path),
		},
		s3.WithPresignExpires(expires),
	)
	if err != nil {
		return nil, err
	}
	return &presignURL.URL, nil
}
