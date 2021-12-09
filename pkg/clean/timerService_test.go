package clean

import (
	"context"
	"testing"
	"time"

	"github.com/airenas/async-api/internal/pkg/test/mocks"
	"github.com/petergtz/pegomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
			IDsProvider: mocks.NewMockOldIDsProvider()}}, wantErr: false},
		{name: "Fail", args: args{ctx: context.Background(), data: &TimerData{RunEvery: time.Second, Cleaner: newCleanMock(false),
			IDsProvider: mocks.NewMockOldIDsProvider()}}, wantErr: true},
		{name: "Fail", args: args{ctx: context.Background(), data: &TimerData{RunEvery: time.Minute, Cleaner: nil,
			IDsProvider: mocks.NewMockOldIDsProvider()}}, wantErr: true},
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
		IDsProvider: mocks.NewMockOldIDsProvider()}
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
		cleanReturn []pegomock.ReturnValue
		idsReturn   []pegomock.ReturnValue
	}
	tests := []struct {
		name          string
		args          args
		cleanExpected pegomock.InvocationCountMatcher
		idsExpected   pegomock.InvocationCountMatcher
	}{
		{name: "Loop", args: args{cleanReturn: []pegomock.ReturnValue{nil},
			idsReturn: []pegomock.ReturnValue{[]string{"1"}, nil}},
			cleanExpected: pegomock.AtLeast(4), idsExpected: pegomock.AtLeast(4)},
		{name: "Several", args: args{cleanReturn: []pegomock.ReturnValue{nil},
			idsReturn: []pegomock.ReturnValue{[]string{"1", "2"}, nil}},
			cleanExpected: pegomock.AtLeast(8), idsExpected: pegomock.AtLeast(4)},
		{name: "No Clean", args: args{cleanReturn: []pegomock.ReturnValue{nil},
			idsReturn: []pegomock.ReturnValue{[]string{}, nil}},
			cleanExpected: pegomock.Never(), idsExpected: pegomock.AtLeast(4)},
		{name: "No Clean - error", args: args{cleanReturn: []pegomock.ReturnValue{nil},
			idsReturn: []pegomock.ReturnValue{nil, errors.New("olia")}},
			cleanExpected: pegomock.Never(), idsExpected: pegomock.AtLeast(4)},
		{name: "Clean - error", args: args{cleanReturn: []pegomock.ReturnValue{errors.New("olia")},
			idsReturn: []pegomock.ReturnValue{[]string{"1"}, nil}},
			cleanExpected: pegomock.AtLeast(4), idsExpected: pegomock.AtLeast(4)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks.AttachMockToTest(t)
			clMock := mocks.NewMockCleaner()
			idsMock := mocks.NewMockOldIDsProvider()
			pegomock.When(idsMock.GetExpired()).ThenReturn(tt.args.idsReturn...)
			pegomock.When(clMock.Clean(pegomock.AnyString())).ThenReturn(tt.args.cleanReturn...)
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
			idsMock.VerifyWasCalled(tt.idsExpected).GetExpired()
			clMock.VerifyWasCalled(tt.cleanExpected).Clean(pegomock.AnyString())
		})
	}
}
