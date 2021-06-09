package flusher_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-snippet-api/internal/flusher"
	"github.com/ozoncp/ocp-snippet-api/internal/mocks"
	"github.com/ozoncp/ocp-snippet-api/internal/models"
)

var _ = Describe("Flush", func() {

	var (
		ctrl *gomock.Controller

		mockRepo      *mocks.MockRepo
		mockPublisher *mocks.MockPublisher

		batchSize uint
		input     []models.Snippet
		rest      []models.Snippet
		err       error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockRepo = mocks.NewMockRepo(ctrl)
		mockPublisher = mocks.NewMockPublisher(ctrl)

		batchSize = 2
		input = []models.Snippet{{Id: 1}, {Id: 2}, {Id: 3}}
	})

	JustBeforeEach(func() {
		fl := flusher.New(batchSize, mockRepo, mockPublisher)
		rest, err = fl.Flush(input)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("repo saves all snippets", func() {

		BeforeEach(func() {
			mockRepo.EXPECT().AddSnippets(gomock.Any()).MinTimes(2).Return(nil)
			mockPublisher.EXPECT().PublishFlushing(gomock.Any())
		})

		It("all saved", func() {
			Expect(len(rest)).To(Equal(0))
			Expect(err).To(BeNil())
		})

	})

	Context("zero batch size", func() {

		BeforeEach(func() {
			batchSize = 0

			mockPublisher.EXPECT().PublishFlushing(0)
		})

		It("nothing to save", func() {
			Expect(len(rest)).To(Equal(len(input)))
			Expect(err).NotTo(BeNil())
		})

	})

	Context("empty input slice", func() {

		BeforeEach(func() {
			input = []models.Snippet{}

			mockRepo.EXPECT().AddSnippets(gomock.Any()).Return(nil)
			mockPublisher.EXPECT().PublishFlushing(0)
		})

		It("nothing to save", func() {
			Expect(len(rest)).To(Equal(0))
			Expect(err).To(BeNil())
		})

	})

	Context("nil input slice", func() {

		BeforeEach(func() {
			input = nil

			mockRepo.EXPECT().AddSnippets(gomock.Any()).Return(nil)
			mockPublisher.EXPECT().PublishFlushing(0)
		})

		It("all saved", func() {
			// проверка рез-та (rest)
			Expect(len(rest)).To(Equal(0))
			Expect(err).To(BeNil())
		})

	})
})
