package quteos

import (
	"app/internal/domain/models"
	"app/internal/storage"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/brianvoe/gofakeit"
)

var ctxB context.Context

var MockStorage map[int]models.Quote = map[int]models.Quote{
	1: {
		Quote:  "Some quote",
		Author: gofakeit.Name(),
	},
	2: {
		Quote:  "bla bla bla bla",
		Author: gofakeit.Name(),
	},
	5: {
		Quote:  "bla bla bla bla",
		Author: gofakeit.Name(),
	},
}

type mockStorage struct {
	storage.Storage
	data map[int]models.Quote
}

func (m *mockStorage) Delete(ctx context.Context, id int) error {
	if _, exists := m.data[id]; !exists {
		return storage.ErrQuoteNotFound
	}
	delete(m.data, id)
	return nil
}

func (m *mockStorage) Get(ctx context.Context, id int) (*models.Quote, error) {

	res, exists := m.data[id]
	if !exists {
		return nil, storage.ErrQuoteNotFound
	}

	return &res, storage.ErrQuoteNotFound
}

func TestService_Save(t *testing.T) {
	type fields struct {
		storage storage.Storage
		log     *slog.Logger
	}
	type args struct {
		ctx context.Context
		q   *models.Quote
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "nil quote",
			fields: fields{
				storage: &mockStorage{},
				log:     slog.Default(),
			},
			args: args{
				ctx: context.Background(),
				q:   nil,
			},
			wantErr: true,
			err:     ErrQuoteIsNil,
		},

		{
			name: "invalid quote (epmty)",
			fields: fields{
				storage: &mockStorage{},
				log:     slog.Default(),
			},
			args: args{
				ctx: context.Background(),
				q: &models.Quote{
					Quote:  "",
					Author: "Test Author",
				},
			},
			wantErr: true,
			err:     ErrValidateQuote,
		},

		{
			name: "invalid author (empty)",
			fields: fields{
				storage: &mockStorage{},
				log:     slog.Default(),
			},
			args: args{
				ctx: context.Background(),
				q: &models.Quote{
					Quote:  "This is a test quote",
					Author: "",
				},
			},
			wantErr: true,
			err:     ErrValidateQuote,
		},

		{
			name: "invalid quote and author < 3 characters",
			fields: fields{
				storage: &mockStorage{},
				log:     slog.Default(),
			},
			args: args{
				ctx: context.Background(),
				q: &models.Quote{
					Quote:  "T",
					Author: "A",
				},
			},
			wantErr: true,
			err:     ErrValidateQuote,
		},

		{
			name: "invalid  author > 500 characters (501 characters)",
			fields: fields{
				storage: &mockStorage{},
				log:     slog.Default(),
			},
			args: args{
				ctx: context.Background(),
				q: &models.Quote{
					Quote:  "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus",
					Author: "Author",
				},
			},
			wantErr: true,
			err:     ErrValidateQuote,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
				log:     tt.fields.log,
			}

			_, err := s.Save(tt.args.ctx, tt.args.q)
			if err != nil {
				if tt.wantErr {
					if !errors.Is(err, tt.err) {
						t.Errorf("Service.Save() error = %v, wantErr %v, expected error %v", err, tt.wantErr, tt.err)
					}
				} else {
					t.Errorf("Service.Save() unexpected error = %v", err)
				}
			}

		})
	}
}

func TestService_Delete(t *testing.T) {

	ctxB = context.Background()

	type fields struct {
		storage storage.Storage
		log     *slog.Logger
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "delete existing quote",
			fields: fields{
				storage: &mockStorage{
					data: MockStorage,
				},
				log: slog.Default(),
			},
			args: args{
				ctx: ctxB,
				id:  "1",
			},
			wantErr: false,
		},
		{
			name: "delete not existing quote",
			fields: fields{
				storage: &mockStorage{
					data: MockStorage,
				},
				log: slog.Default(),
			},
			args: args{
				ctx: ctxB,
				id:  "3",
			},
			wantErr: true,
			err:     storage.ErrQuoteNotFound,
		},
		{
			name: "not valid quote id",
			fields: fields{
				storage: &mockStorage{
					data: MockStorage,
				},
				log: slog.Default(),
			},
			args: args{
				ctx: ctxB,
				id:  "abc",
			},
			wantErr: true,
			err:     ErrInvalidQuoteID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
				log:     tt.fields.log,
			}

			err := s.Delete(tt.args.ctx, tt.args.id)
			if err != nil {
				if tt.wantErr {
					if !errors.Is(err, tt.err) {
						t.Errorf("Service.Delete() error = %v, wantErr %v, expected error %v", err, tt.wantErr, tt.err)
					}

				}

			}

			_, err = s.Get(ctxB, tt.args.id)
			if err != nil {
				if tt.wantErr {
					if errors.Is(err, storage.ErrQuoteNotFound) {
						t.Logf("Quote with id %s successfully deleted", tt.args.id)
					}

					if !errors.Is(err, tt.err) {
						t.Errorf("Service.Get() error = %v, wantErr %v, expected error %v", err, tt.wantErr, tt.err)
					}
				}
			}
		})
	}
}
