# Deprecated

**DEPRECATED:** Use
**[`picocolors`](https://github.com/alexeyraspopov/picocolors)** instead.
It is 3 times smaller and 50% faster.

The space in node_modules including sub-dependencies:

```diff
- nanocolors   16 kB
+ picocolors    7 kB
```

Library loading time:

```diff
- nanocolors     0.885 ms
+ picocolors     0.470 ms
```

Benchmark for complex use cases:

```diff
- nanocolors     1,088,193 ops/sec
+ picocolors     1,772,265 ops/sec
```

## Docs
Read **[full docs](https://github.com/ai/nanocolors#readme)** on GitHub.
