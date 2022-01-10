## Remaining EE Codebase

### Folders
- ``Analysis`` includes tools and data (copied from Google sheets) used to produce the graphs included in the paper
- ``notes`` includes various notes used during the development process, about the buffers and how to test them
- ``old`` includes some code from the first iteration of testing, and just some miscellaneous stuff that isn't too relevant
- ``results`` includes the ``.txt`` files that include the data copied off of the benchmarking results screen, as well as the ``.csv`` files generated by running ``results_reader.go``
- ``rope-vis`` includes the input file and image geneated by the Graphviz program after using the associated Golang bindings, ``memviz``: https://github.com/bradleyjkemp/memviz. The associated methods are included in ``notes`` under ``misc_test.md``

### Other files
The remainder of the stuff is just the actual data structures and tests for them, the results reader, tree utilities for testing the rope (``tree_utils.go``), and the two text files used for initial and final testing (``testing.txt`` and ``pg66576.txt`` respectively.)