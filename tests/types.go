//go:build tangentgen

package tests

import tangent "github.com/telophasehq/tangent-sdk-go"

// _gen only exists so the generator can find Wire[T] instantiations.
// It's never called.
func _gen() {
	// Mention every struct you want generated:
	tangent.Wire[MyStructAnon](tangent.Metadata{}, nil, nil)
	tangent.Wire[MyStructAnonPtr](tangent.Metadata{}, nil, nil)
	tangent.Wire[MyStruct](tangent.Metadata{}, nil, nil)
	tangent.Wire[MyStructPtr](tangent.Metadata{}, nil, nil)
	tangent.Wire[MyNested](tangent.Metadata{}, nil, nil)
}
