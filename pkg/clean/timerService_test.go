package clean

import (
	"context"
	"testing"
	"time"

	"github.com/airenas/async-api/internal/pkg/test"
	"github.com/airenas/async-api/internal/pkg/test/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartCleanTimer(t *testing.T) {
	type args struct {
		ctx  context.Context
		data *TimerData
	}
	tests := []struct {
		name    string
		args    args
		want    <-chan struct{}
		wantErr bool
	}{
		{name: "Init", args: args{ctx: context.Background(), data: &TimerData{RunEvery: time.Hour, Cleaner: newCleanMock(false),
			IDsProvider: newIDsProviderMock(nil, false)}}, wantErr: false},
		{name: "Fail", args: args{ctx: context.Background(), data: &TimerData{RunEvery: time.Second, Cleaner: newCleanMock(false),
			IDsProvider: newIDsProviderMock(nil, false)}}, wantErr: true},
		{name: "Fail", args: args{ctx: context.Background(), data: &TimerData{RunEvery: time.Minute, Cleaner: nil,
			IDsProvider: newIDsProviderMock(nil, false)}}, wantErr: true},
		{name: "Fail", args: args{ctx: context.Background(), data: &TimerData{RunEvery: time.Minute, Cleaner: newCleanMock(false),
			IDsProvider: nil}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StartCleanTimer(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartCleanTimer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestExitLoop(t *testing.T) {
	data := &TimerData{RunEvery: time.Second, Cleaner: newCleanMock(false),
		IDsProvider: newIDsProviderMock(nil, false)}
	ctx, cFunc := context.WithCancel(context.Background())
	ch := startLoop(ctx, data)
	cFunc()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Errorf("Timeout")
	}
}

func TestTimerLoops(t *testing.T) {
	type args struct {
		cleanErr  bool
		idsReturn []string
		idsErr    bool
	}
	tests := []struct {
		name          string
		args          args
		cleanExpected int
		idsExpected   int
	}{
		{name: "Loop", args: args{cleanErr: false,
			idsReturn: []string{"1"}, idsErr: false},
			cleanExpected: 4, idsExpected: 4},
		{name: "Several", args: args{cleanErr: false,
			idsReturn: []string{"1", "2"}, idsErr: false},
			cleanExpected: 8, idsExpected: 4},
		{name: "No Clean - error", args: args{cleanErr: false,
			idsReturn: nil, idsErr: true},
			cleanExpected: 0, idsExpected: 4},
		{name: "Clean - error", args: args{cleanErr: true,
			idsReturn: []string{"1"}, idsErr: false},
			cleanExpected: 4, idsExpected: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks.AttachMockToTest(t)
			clMock := newCleanMock(tt.args.cleanErr)
			idsMock := newIDsProviderMock(tt.args.idsReturn, tt.args.idsErr)
			data := &TimerData{RunEvery: time.Millisecond * 5, Cleaner: clMock,
				IDsProvider: idsMock}
			ctx, cFunc := context.WithCancel(context.Background())
			ch := startLoop(ctx, data)
			<-time.After(time.Millisecond * 20)
			cFunc()
			select {
			case <-ch:
			case <-time.After(time.Second):
				t.Errorf("Timeout")
			}
			assert.LessOrEqual(t, tt.cleanExpected, len(clMock.Calls))
			assert.LessOrEqual(t, tt.idsExpected, len(idsMock.Calls))
		})
	}
}

type mockIDsProvider struct{ mock.Mock }

func (m *mockIDsProvider) GetExpired(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return test.To[[]string](args.Get(0)), args.Error(1)
}

func newIDsProviderMock(resIDs []string, fail bool) *mockIDsProvider {
	res := &mockIDsProvider{}
	var err error
	if fail {
		err = errors.New("olia")
	}
	res.On("GetExpired", mock.Anything).Return(resIDs, err)
	return res
}
