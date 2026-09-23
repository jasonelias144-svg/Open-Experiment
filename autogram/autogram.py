"""Search for autograms: sentences that correctly inventory their own letters.

Example shape: "This sentence contains three a's, one b, ... and two z's."

Method (after Lee Sallows): treat the sentence as a function f(counts) -> counts.
An autogram is a fixed point. Plain iteration x -> f(x) tends to fall into cycles,
so we move only a random subset of letters toward f(x) each step, and restart
when stuck.
"""
import random
import string
import sys
from collections import Counter

ONES = ["zero", "one", "two", "three", "four", "five", "six", "seven", "eight",
        "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen",
        "sixteen", "seventeen", "eighteen", "nineteen"]
TENS = ["", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy",
        "eighty", "ninety"]

def number_word(n):
    if n < 20:
        return ONES[n]
    t, o = divmod(n, 10)
    return TENS[t] + ("-" + ONES[o] if o else "")

def render(prefix, counts, letters):
    parts = []
    for c in letters:
        n = counts[c]
        parts.append(f"{number_word(n)} {c}" + ("" if n == 1 else "'s"))
    return prefix + ", ".join(parts[:-1]) + ", and " + parts[-1] + "."

def tally(text, letters):
    got = Counter(ch for ch in text.lower() if ch in letters)
    return {c: got[c] for c in letters}

def search(prefix, letters=string.ascii_lowercase, max_steps=2_000_000,
           seed=None, restart_every=20_000):
    rng = random.Random(seed)
    base = tally(prefix, letters)
    counts = {c: base[c] + rng.randint(1, 10) for c in letters}
    for step in range(max_steps):
        actual = tally(render(prefix, counts, letters), letters)
        if actual == counts:
            return step, counts
        if step % restart_every == restart_every - 1:
            counts = {c: max(1, actual[c] + rng.randint(-3, 3)) for c in letters}
            continue
        for c in letters:
            if actual[c] != counts[c] and rng.random() < 0.3:
                counts[c] = actual[c]
    return None

if __name__ == "__main__":
    prefix = sys.argv[1] if len(sys.argv) > 1 else "This sentence contains "
    seed = int(sys.argv[2]) if len(sys.argv) > 2 else None
    result = search(prefix, seed=seed)
    if result is None:
        print("no fixed point found")
        sys.exit(1)
    step, counts = result
    s = render(prefix, counts, string.ascii_lowercase)
    assert tally(s, string.ascii_lowercase) == counts
    print(f"found after {step} steps:\n\n{s}")
