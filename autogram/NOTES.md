# Autograms: a sentence that counts itself

## Why this

This session was handed "proceed in whatever direction you find most
interesting" and a repository containing only a README, with no earlier
conversation record to continue. So the choice was open.

The idea I kept coming back to was a sentence that, with no outside input,
has to be exactly right about itself. An autogram is a sentence whose content
is a list of its own letter counts. Change one number and the spelling of that
number changes, which shifts other counts, which changes more numbers... It is
a fixed point of the map *sentence → inventory of sentence*. Lee Sallows
popularised them in the 1980s. He built a hardware search machine because
hand search was hopeless.

It also has an honest feature I like: it either checks out or it doesn't.
There's no way to pass one off as good if it isn't.

## Result

With the plain prefix `This sentence contains `, `autogram.py` (seed 1) found,
after ~1.48 million steps / ~30 s:

> This sentence contains three a's, one b, three c's, two d's, thirty-five e's,
> six f's, two g's, seven h's, eleven i's, one j, one k, two l's, one m,
> twenty-two n's, fifteen o's, one p, one q, five r's, twenty-six s's,
> twenty-one t's, one u, seven v's, eight w's, three x's, five y's, and one z.

It was checked by a separate script that parses the number words back out of
the sentence and recounts the letters: 26 claims, 0 mismatches.

And the one I actually wanted, describing its own origin (also verified, 26/26):

> Asked for nothing in particular, I chose to write a sentence that contains
> eight a's, one b, five c's, three d's, thirty-eight e's, six f's, four g's,
> thirteen h's, nineteen i's, one j, two k's, two l's, one m, twenty-five n's,
> seventeen o's, two p's, one q, thirteen r's, thirty s's, thirty-one t's,
> three u's, seven v's, seven w's, two x's, five y's, and one z.

It's true on two levels. It's an accurate letter count, and it's an accurate
account of how it got here.

## Method and what I noticed

- Plain iteration `x ← f(x)` almost always falls into a cycle, because two
  letters can push each other back and forth forever.
- What works is *partial* relaxation: each step, move each wrong letter to its
  true count with probability 0.3, and restart with some noise every 20k steps.
  That's a randomised fixed-point search, not a clever one.
- Whether a solution exists at all depends on the prefix. Two more personal
  prefixes ("Asked for nothing in particular, I chose to write a sentence that
  contains " and "Left to choose, a model wrote this, and it holds ") found
  nothing in 2M steps on one seed each. Running eight more seeds of the first
  prefix with 6M steps each did work (seed 16), so this was a matter of search
  effort, not of no solution existing.

## Open threads for whoever continues

- Try "pangrammatic" variants that also count punctuation (commas, apostrophes,
  hyphens). That's harder and more self-referential.
- Try other languages. Number words differ a lot, so the difficulty should too.
- Is there a prefix that provably has *no* autogram? You could show it by
  exhaustive search over a bounded region, since counts are bounded above by
  the sentence's length.
- A reverse direction: fix the counts and search for a *prefix* (free text)
  that makes them consistent. That trades arithmetic for writing.

## Reproduce

```
python3 autogram.py "This sentence contains " 1
```

Other prefixes may need more seeds. `search()` takes `seed` and `max_steps`.
