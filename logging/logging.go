/**
 * @file    logging.go
 * @author  Christophe Buffard
 * @date    2023-06-17
 * @brief   DNSHelper Service for CoreDNS.
 *
 * License under GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007
 * service Logging
 */

package logging

import (
	"context"
	"fmt"
	"github.com/Go-routine-4995/ossfrontend/domain"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type IService interface {
	CreateRouters(ctx context.Context, r []domain.Router, tenant string) ([]byte, error)
	DeleteRouters(ctx context.Context, r []domain.Router, tenant string) error
	GetRoutersPage(ctx context.Context, paginationByte []byte, tenant string) ([]domain.Router, error)
	GetRouters(ctx context.Context, r domain.Router, tenant string) (domain.Router, error)
}

type LoggingService struct {
	next IService
	log  zerolog.Logger
}

func NewLoggingService(next interface{}) IService {

	return &LoggingService{
		next: next.(IService),
		log:  zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.InfoLevel).With().Timestamp().Logger(),
	}

}

func (s *LoggingService) CreateRouters(ctx context.Context, r []domain.Router, tenant string) (rep []byte, err error) {

	defer func(start time.Time) {
		var str string
		var sreq string
		if rep != nil {
			str = fmt.Sprintf("%v", rep)
		}
		if len(r) < 5 {
			sreq = fmt.Sprintf("%v", r)
		} else {
			sreq = fmt.Sprintf("request too large: %d routers being created", len(r))
		}
		s.log.Info().
			Str("method", "CreateRouters").
			Str("request", sreq).
			Str("response", str).
			Str("tenant", tenant).
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.CreateRouters(ctx, r, tenant)
}

func (s *LoggingService) DeleteRouters(ctx context.Context, r []domain.Router, tenant string) (err error) {

	defer func(start time.Time) {
		var sreq string
		if len(r) < 5 {
			sreq = fmt.Sprintf("%v", r)
		} else {
			sreq = fmt.Sprintf("request too large: %d routers being created", len(r))
		}
		s.log.Info().
			Str("method", "DeleteRouters").
			Str("request", sreq).
			Str("tenant", tenant).
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.DeleteRouters(ctx, r, tenant)
}

func (s *LoggingService) GetRoutersPage(ctx context.Context, paginationByte []byte, tenant string) (rep []domain.Router, err error) {

	defer func(start time.Time) {
		var str string
		var sreq string
		if rep != nil {
			if len(rep) < 5 {
				str = fmt.Sprintf("%v", rep)
			} else {
				str = fmt.Sprintf("%v", rep)
			}
		}

		sreq = fmt.Sprintf("%v", paginationByte)

		s.log.Info().
			Str("method", "GetRoutersPaged").
			Str("request", sreq).
			Str("response", str).
			Str("tenant", tenant).
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.GetRoutersPage(ctx, paginationByte, tenant)
}

func (s *LoggingService) GetRouters(ctx context.Context, r domain.Router, tenant string) (rep domain.Router, err error) {

	defer func(start time.Time) {
		var str string
		var sreq string

		str = fmt.Sprintf("%v", rep)

		sreq = fmt.Sprintf("%v", r)

		s.log.Info().
			Str("method", "GetRouter").
			Str("request", sreq).
			Str("response", str).
			Str("tenant", tenant).
			Err(err).
			Dur("took", time.Since(start)).Send()
	}(time.Now())

	return s.next.GetRouters(ctx, r, tenant)
}
