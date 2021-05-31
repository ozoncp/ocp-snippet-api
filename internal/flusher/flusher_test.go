package flusher_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	"github.com/ozoncp/ocp-snippet-api/internal/mocks"
)

var _ = Describe("Flusher", func() {

	var (
		ctrl *gomock.Controller

		mockRepo      *mocks.MockRepo
		mockPublisher *mocks.MockPublisher
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockRepo = mocks.NewMockRepo(ctrl)
		mockPublisher = mocks.NewMockPublisher(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("repo saves all tasks", func() {

		BeforeEach(func() {
			mockRepo.EXPECT().AddTasks(gomock.Any()).Return(nil)
			mockPublisher.EXPECT().PublishFlushing(gomock.Any())
		})
	})
})
