package stylegen

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestParseRecipesValid(t *testing.T) {
	t.Parallel()

	filename := filepath.Join("testdata", "recipe-valid.css")
	src, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	recipes, err := ParseRecipes(filename, src)
	if err != nil {
		t.Fatalf("ParseRecipes() error = %v", err)
	}

	if got, want := recipes.Tokens(), []string{
		"gsxui-recipe-button",
		"gsxui-recipe-button-variant-outline",
		"gsxui-recipe-button-size-icon",
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("Tokens() = %q, want %q", got, want)
	}

	recipe, ok := recipes.Lookup("gsxui-recipe-button")
	if !ok {
		t.Fatal("Lookup(gsxui-recipe-button) = false, want true")
	}
	if got, want := recipe.Utilities, []string{
		"inline-flex",
		"items-center",
		"gap-2",
		"hover:bg-primary/90",
		"data-\\[state=open\\]:bg-primary",
		"[&_svg]:size-4",
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("Lookup(gsxui-recipe-button).Utilities = %q, want %q", got, want)
	}

	if _, ok := recipes.Lookup("gsxui-recipe-missing"); ok {
		t.Error("Lookup(gsxui-recipe-missing) = true, want false")
	}
}

func TestParseRecipesRejectsInvalidGrammar(t *testing.T) {
	t.Parallel()

	const recipe = ".gsxui-recipe-button"
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "no components layer",
			src:  recipe + " { @apply flex; }",
			want: recipe,
		},
		{
			name: "wrong layer",
			src: `@layer utilities {
  .gsxui-recipe-button { @apply flex; }
}`,
			want: "@layer utilities",
		},
		{
			name: "nested layer",
			src: `@layer components {
  @layer utilities {
    .gsxui-recipe-button { @apply flex; }
  }
}`,
			want: "@layer utilities",
		},
		{
			name: "non recipe rule",
			src: `@layer components {
  .button { @apply flex; }
}`,
			want: ".button",
		},
		{
			name: "empty rule",
			src: `@layer components {
  .gsxui-recipe-button {}
}`,
			want: recipe,
		},
		{
			name: "ordinary declaration",
			src: `@layer components {
  .gsxui-recipe-button { color: red; }
}`,
			want: recipe,
		},
		{
			name: "duplicate recipe token",
			src: `@layer components {
  .gsxui-recipe-button { @apply flex; }
  .gsxui-recipe-button { @apply grid; }
}`,
			want: recipe,
		},
		{
			name: "element selector",
			src: `@layer components {
  button.gsxui-recipe-button { @apply flex; }
}`,
			want: "button.gsxui-recipe-button",
		},
		{
			name: "id selector",
			src: `@layer components {
  #gsxui-recipe-button { @apply flex; }
}`,
			want: "#gsxui-recipe-button",
		},
		{
			name: "attribute selector",
			src: `@layer components {
  .gsxui-recipe-button[data-size] { @apply flex; }
}`,
			want: ".gsxui-recipe-button[data-size]",
		},
		{
			name: "descendant selector",
			src: `@layer components {
  .gsxui-recipe-button .icon { @apply flex; }
}`,
			want: ".gsxui-recipe-button .icon",
		},
		{
			name: "child selector",
			src: `@layer components {
  .gsxui-recipe-button > .icon { @apply flex; }
}`,
			want: ".gsxui-recipe-button>.icon",
		},
		{
			name: "adjacent sibling selector",
			src: `@layer components {
  .gsxui-recipe-button + .icon { @apply flex; }
}`,
			want: ".gsxui-recipe-button+.icon",
		},
		{
			name: "general sibling selector",
			src: `@layer components {
  .gsxui-recipe-button ~ .icon { @apply flex; }
}`,
			want: ".gsxui-recipe-button~.icon",
		},
		{
			name: "pseudo selector",
			src: `@layer components {
  .gsxui-recipe-button:hover { @apply flex; }
}`,
			want: ".gsxui-recipe-button:hover",
		},
		{
			name: "compound selector",
			src: `@layer components {
  .gsxui-recipe-button.primary { @apply flex; }
}`,
			want: ".gsxui-recipe-button.primary",
		},
		{
			name: "comma selector list",
			src: `@layer components {
  .gsxui-recipe-button, .other { @apply flex; }
}`,
			want: ".gsxui-recipe-button,.other",
		},
		{
			name: "escaped selector name",
			src: `@layer components {
  .gsxui-\72 ecipe-button { @apply flex; }
}`,
			want: `.gsxui-\72 ecipe-button`,
		},
		{
			name: "constructed selector name",
			src: `@layer components {
  .gsxui-recipe-button#{suffix} { @apply flex; }
}`,
			want: "gsxui-recipe-button",
		},
		{
			name: "nested rule",
			src: `@layer components {
  .gsxui-recipe-button {
    @apply flex;
    .icon { @apply size-4; }
  }
}`,
			want: ".icon",
		},
		{
			name: "empty apply",
			src: `@layer components {
  .gsxui-recipe-button { @apply; }
}`,
			want: "@apply",
		},
		{
			name: "unterminated utility",
			src: `@layer components {
  .gsxui-recipe-button { @apply hover:; }
}`,
			want: "@apply",
		},
		{
			name: "malformed css",
			src: `@layer components {
  .gsxui-recipe-button { @apply flex; `,
			want: recipe,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join("testdata", strings.ReplaceAll(tt.name, " ", "-")+".css")
			_, err := ParseRecipes(filename, []byte(tt.src))
			if err == nil {
				t.Fatal("ParseRecipes() error = nil, want error")
			}

			message := err.Error()
			if !strings.Contains(message, filename) {
				t.Errorf("error %q does not include filename %q", message, filename)
			}
			if !regexp.MustCompile(regexp.QuoteMeta(filename) + `:\d+:\d+:`).MatchString(message) {
				t.Errorf("error %q does not include filename, line, and column", message)
			}
			if !strings.Contains(message, tt.want) {
				t.Errorf("error %q does not include offending %q", message, tt.want)
			}
		})
	}
}
