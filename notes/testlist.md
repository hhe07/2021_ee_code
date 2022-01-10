New Testing Standards
- Load (done)
  - Test with a few variations of load input size
- Inserts (done)
  - At start
  - At end
  - At centre
- Reports (done)
  - Test with a few variations of load input size
- Report Ranges (done)
  - Test with a few variations of ranges:
    - Whole
    - One side at start
    - One side at end
    - At a few intermediate intervals
- Appends (done)
  - Test with a few variations of input size
- Split (done)
  - Test with a few variations of split point at middle
- DeleteRange
  - One side at start
  - One side at end
  - At a few intermediate ranges (similar to report)


New Benchmarking Standards

first, see if some sort of propagatory delay works

- Load (done)
  - Start: 50K
  - Step: 50K
  - Range: 100K - 120K
- Report (same as previous) (done)
- Report Range (done)
  - First, test whether size of buffer has any influence on report times at a constant position:
    - (same as previous)
    - report 10k characters at an arbitrary position
  - Second, test whether length of report has any influence on report times
    - (constant buffer size)
  - TODO: these had poor results
- Insert (done)
  - Buffer size?
  - First, scale input size
  - Second, scale position
  - TODO: rope ones were bad
- Split (done)
  - Buffer size?
  - Scale position
  - TODO: both bad
- Delete
  - Buffer size?
  - Length
  - Position
  - TODO: rope bad










