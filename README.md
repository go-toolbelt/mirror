# mirror
A small reflection utility for Go for high-performance java-style stack traces

Designed to have no uncached heap allocations. The trade-off is that stack traces are limited to a fixed depth (32 frames).

Formatted stack frame values are computed once and cached. This prevents additional allocations for any subsequent time that a stack frame it taken for any code point.
