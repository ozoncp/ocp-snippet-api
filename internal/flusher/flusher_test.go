package flusher_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

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
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockRepo = mocks.NewMockRepo(ctrl)
		mockPublisher = mocks.NewMockPublisher(ctrl)

		batchSize = 2
		input = []models.Snippet{{Id: 1}, {Id: 2}, {Id: 3}}
	})

	JustBeforeEach(func() {
		// Создать экземпляр Flush с помощью New( mockRepo )
		// позвать Flush
		// записать рез-т в rest

		fl := flusher.New(batchSize, mockRepo, mockPublisher)
		rest = fl.Flush(input)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("repo saves all snippets", func() {

		BeforeEach(func() {
			mockRepo.EXPECT().AddSnippets(gomock.Any()).Return(nil)
			mockPublisher.EXPECT().PublishFlushing(gomock.Any())
		})

		It("all saved", func() {
			// проверка рез-та (rest)
			gomega.Expect(len(rest)).To(gomega.Equal(0))
		})

	})

	Context("zero batch size", func() {

		BeforeEach(func() {
			batchSize = 0

			mockPublisher.EXPECT().PublishFlushing(0)
		})

		It("nothing to save", func() {
			// проверка рез-та (rest)
			gomega.Expect(len(rest)).To(gomega.Equal(len(input)))
		})

	})

	Context("empty input slice", func() {

		BeforeEach(func() {
			input = []models.Snippet{}

			mockRepo.EXPECT().AddSnippets(gomock.Any()).Return(nil)
			mockPublisher.EXPECT().PublishFlushing(0)
		})

		It("nothing to save", func() {
			// проверка рез-та (rest)
			gomega.Expect(len(rest)).To(gomega.Equal(0))
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
			gomega.Expect(len(rest)).To(gomega.Equal(0))
		})

	})
})
