package runner_test

import (
	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/cucumber/cucumber-engine/src/runner"
	gherkin "github.com/cucumber/gherkin-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PickleFilter", func() {
	Describe("Matches", func() {
		var pickleFilter *runner.PickleFilter
		var pickleEvent *gherkin.PickleEvent

		BeforeEach(func() {
			pickleEvent = &gherkin.PickleEvent{
				Pickle: &gherkin.Pickle{
					Name:      "",
					Tags:      []*gherkin.PickleTag{},
					Locations: []gherkin.Location{},
				},
				URI: "",
			}
		})

		Describe("no filters", func() {
			BeforeEach(func() {
				var err error
				pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
					Lines:         map[string][]int{},
					Names:         []string{},
					TagExpression: "",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns true", func() {
				Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
			})
		})

		Describe("line filters", func() {
			BeforeEach(func() {
				var err error
				pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
					Lines:         map[string][]int{"/path/to/featureB": []int{1, 2}},
					Names:         []string{},
					TagExpression: "",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			Describe("scenario in feature without line specified", func() {
				BeforeEach(func() {
					pickleEvent.URI = "/path/to/featureA"
				})

				It("returns true", func() {
					Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
				})
			})

			Describe("scenario in feature with line specified", func() {
				BeforeEach(func() {
					pickleEvent.URI = "/path/to/featureB"
				})

				Describe("first pickle line matches", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Locations = []gherkin.Location{{Line: 1}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle line does not match", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Locations = []gherkin.Location{{Line: 3}}
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})
			})
		})

		Describe("name filters", func() {
			Describe("one name", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
						Names:         []string{"nameA"},
						TagExpression: "",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle name includes filter", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Name = "nameA descriptionA"
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle name does not include filter", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Name = "nameB descriptionB"
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})
			})

			Describe("multiple names", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
						Names:         []string{"nameA", "nameB"},
						TagExpression: "",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle name includes the first filter", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Name = "nameA descriptionA"
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle name includes the second filter", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Name = "nameB descriptionB"
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle name does not include either filter", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Name = "nameC descriptionC"
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})
			})
		})

		Describe("tag filters", func() {
			Describe("filtering with a single tag", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
						Lines:         map[string][]int{},
						Names:         []string{},
						TagExpression: "@tagA",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle has tag", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagA"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle does not have tag", func() {
					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})
			})

			Describe("filtering with a negation of a single tag", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
						Lines:         map[string][]int{},
						Names:         []string{},
						TagExpression: "not @tagA",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle has tag", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagA"}}
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})

				Describe("pickle does not have tag", func() {
					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})
			})

			Describe("filtering by and-ing two tags", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
						Lines:         map[string][]int{},
						Names:         []string{},
						TagExpression: "@tagA and @tagB",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle has both tags", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagA"}, {Name: "@tagB"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle has the first tag but not the second", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagA"}}
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})

				Describe("pickle has the second tag but not the first", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagB"}}
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})

				Describe("pickle does not either tag", func() {
					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})
			})

			Describe("filtering by or-ing two tags", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
						Lines:         map[string][]int{},
						Names:         []string{},
						TagExpression: "@tagA or @tagB",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle has both tags", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagA"}, {Name: "@tagB"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle has the first tag but not the second", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagA"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle has the second tag but not the first", func() {
					BeforeEach(func() {
						pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagB"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
					})
				})

				Describe("pickle does not either tag", func() {
					It("returns false", func() {
						Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
					})
				})
			})
		})

		Describe("line, name, and tag filters", func() {
			BeforeEach(func() {
				var err error
				pickleFilter, err = runner.NewPickleFilter(&dto.FeaturesFilterConfig{
					Lines:         map[string][]int{"features/b.feature": []int{1, 2}},
					Names:         []string{"nameA"},
					TagExpression: "@tagA",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			Describe("pickle matches all filters", func() {
				BeforeEach(func() {
					pickleEvent.URI = "features/b.feature"
					pickleEvent.Pickle.Locations = []gherkin.Location{{Line: 1}}
					pickleEvent.Pickle.Name = "nameA descriptionA"
					pickleEvent.Pickle.Tags = []*gherkin.PickleTag{{Name: "@tagA"}}
				})

				It("returns true", func() {
					Expect(pickleFilter.Matches(pickleEvent)).To(BeTrue())
				})
			})

			Describe("pickle matches some filters", func() {
				BeforeEach(func() {
					pickleEvent.URI = "features/b.feature"
					pickleEvent.Pickle.Locations = []gherkin.Location{{Line: 1}}
				})

				It("returns false", func() {
					Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
				})
			})

			Describe("pickle matches no filters", func() {
				It("returns false", func() {
					Expect(pickleFilter.Matches(pickleEvent)).To(BeFalse())
				})
			})
		})
	})
})
