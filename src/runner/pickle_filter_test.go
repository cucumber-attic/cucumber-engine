package runner_test

import (
	"github.com/cucumber/cucumber-engine/src/runner"
	messages "github.com/cucumber/cucumber-messages-go/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PickleFilter", func() {
	Describe("Matches", func() {
		var pickleFilter *runner.PickleFilter
		var pickle *messages.Pickle

		BeforeEach(func() {
			pickle = &messages.Pickle{
				Name:      "",
				Tags:      []*messages.PickleTag{},
				Locations: []*messages.Location{},
				Uri:       "",
			}
		})

		Describe("no filters", func() {
			BeforeEach(func() {
				var err error
				pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
					UriToLinesMapping:      []*messages.UriToLinesMapping{},
					NameRegularExpressions: []string{},
					TagExpression:          "",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns true", func() {
				Expect(pickleFilter.Matches(pickle)).To(BeTrue())
			})
		})

		Describe("line filters", func() {
			BeforeEach(func() {
				var err error
				pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
					UriToLinesMapping: []*messages.UriToLinesMapping{
						{
							AbsolutePath: "/path/to/featureB",
							Lines:        []uint64{1, 2},
						},
					},
					NameRegularExpressions: []string{},
					TagExpression:          "",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			Describe("scenario in feature without line specified", func() {
				BeforeEach(func() {
					pickle.Uri = "/path/to/featureA"
				})

				It("returns true", func() {
					Expect(pickleFilter.Matches(pickle)).To(BeTrue())
				})
			})

			Describe("scenario in feature with line specified", func() {
				BeforeEach(func() {
					pickle.Uri = "/path/to/featureB"
				})

				Describe("first pickle line matches", func() {
					BeforeEach(func() {
						pickle.Locations = []*messages.Location{{Line: 1}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle line does not match", func() {
					BeforeEach(func() {
						pickle.Locations = []*messages.Location{{Line: 3}}
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})
			})
		})

		Describe("name filters", func() {
			Describe("one name", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
						UriToLinesMapping:      []*messages.UriToLinesMapping{},
						NameRegularExpressions: []string{"nameA"},
						TagExpression:          "",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle name includes filter", func() {
					BeforeEach(func() {
						pickle.Name = "nameA descriptionA"
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle name does not include filter", func() {
					BeforeEach(func() {
						pickle.Name = "nameB descriptionB"
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})
			})

			Describe("multiple names", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
						UriToLinesMapping:      []*messages.UriToLinesMapping{},
						NameRegularExpressions: []string{"nameA", "nameB"},
						TagExpression:          "",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle name includes the first filter", func() {
					BeforeEach(func() {
						pickle.Name = "nameA descriptionA"
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle name includes the second filter", func() {
					BeforeEach(func() {
						pickle.Name = "nameB descriptionB"
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle name does not include either filter", func() {
					BeforeEach(func() {
						pickle.Name = "nameC descriptionC"
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})
			})
		})

		Describe("tag filters", func() {
			Describe("filtering with a single tag", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
						UriToLinesMapping:      []*messages.UriToLinesMapping{},
						NameRegularExpressions: []string{},
						TagExpression:          "@tagA",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle has tag", func() {
					BeforeEach(func() {
						pickle.Tags = []*messages.PickleTag{{Name: "@tagA"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle does not have tag", func() {
					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})
			})

			Describe("filtering with a negation of a single tag", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
						UriToLinesMapping:      []*messages.UriToLinesMapping{},
						NameRegularExpressions: []string{},
						TagExpression:          "not @tagA",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle has tag", func() {
					BeforeEach(func() {
						pickle.Tags = []*messages.PickleTag{{Name: "@tagA"}}
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})

				Describe("pickle does not have tag", func() {
					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})
			})

			Describe("filtering by and-ing two tags", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
						UriToLinesMapping:      []*messages.UriToLinesMapping{},
						NameRegularExpressions: []string{},
						TagExpression:          "@tagA and @tagB",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle has both tags", func() {
					BeforeEach(func() {
						pickle.Tags = []*messages.PickleTag{{Name: "@tagA"}, {Name: "@tagB"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle has the first tag but not the second", func() {
					BeforeEach(func() {
						pickle.Tags = []*messages.PickleTag{{Name: "@tagA"}}
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})

				Describe("pickle has the second tag but not the first", func() {
					BeforeEach(func() {
						pickle.Tags = []*messages.PickleTag{{Name: "@tagB"}}
					})

					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})

				Describe("pickle does not either tag", func() {
					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})
			})

			Describe("filtering by or-ing two tags", func() {
				BeforeEach(func() {
					var err error
					pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
						UriToLinesMapping:      []*messages.UriToLinesMapping{},
						NameRegularExpressions: []string{},
						TagExpression:          "@tagA or @tagB",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("pickle has both tags", func() {
					BeforeEach(func() {
						pickle.Tags = []*messages.PickleTag{{Name: "@tagA"}, {Name: "@tagB"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle has the first tag but not the second", func() {
					BeforeEach(func() {
						pickle.Tags = []*messages.PickleTag{{Name: "@tagA"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle has the second tag but not the first", func() {
					BeforeEach(func() {
						pickle.Tags = []*messages.PickleTag{{Name: "@tagB"}}
					})

					It("returns true", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeTrue())
					})
				})

				Describe("pickle does not either tag", func() {
					It("returns false", func() {
						Expect(pickleFilter.Matches(pickle)).To(BeFalse())
					})
				})
			})
		})

		Describe("line, name, and tag filters", func() {
			BeforeEach(func() {
				var err error
				pickleFilter, err = runner.NewPickleFilter(&messages.SourcesFilterConfig{
					UriToLinesMapping: []*messages.UriToLinesMapping{
						{
							AbsolutePath: "features/b.feature",
							Lines:        []uint64{1, 2},
						},
					},
					NameRegularExpressions: []string{"nameA"},
					TagExpression:          "@tagA",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			Describe("pickle matches all filters", func() {
				BeforeEach(func() {
					pickle.Uri = "features/b.feature"
					pickle.Locations = []*messages.Location{{Line: 1}}
					pickle.Name = "nameA descriptionA"
					pickle.Tags = []*messages.PickleTag{{Name: "@tagA"}}
				})

				It("returns true", func() {
					Expect(pickleFilter.Matches(pickle)).To(BeTrue())
				})
			})

			Describe("pickle matches some filters", func() {
				BeforeEach(func() {
					pickle.Uri = "features/b.feature"
					pickle.Locations = []*messages.Location{{Line: 1}}
				})

				It("returns false", func() {
					Expect(pickleFilter.Matches(pickle)).To(BeFalse())
				})
			})

			Describe("pickle matches no filters", func() {
				It("returns false", func() {
					Expect(pickleFilter.Matches(pickle)).To(BeFalse())
				})
			})
		})
	})
})
