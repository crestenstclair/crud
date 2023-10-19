package crud

import ( "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/crestenstclair/crud/internal/config"
	"github.com/crestenstclair/crud/internal/repo"
	"github.com/crestenstclair/crud/internal/repo/dynamo"
	"go.uber.org/zap"
)

type Crud struct {
	Repo   repo.Repo
	Logger *zap.Logger
	Config *config.Config
}

func New() (*Crud, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	sess := session.Must(session.NewSession())
	client := dynamodb.New(sess)
	repo, err := dynamo.New(cfg.DYNAMODB_TABLE, client)
	return &Crud{
		Logger: logger,
		Repo:   repo,
		Config: cfg,
	}, nil
}
