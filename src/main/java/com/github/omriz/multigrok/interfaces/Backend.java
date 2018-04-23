package com.github.omriz.multigrok.interfaces;

import com.github.omriz.multigrok.datatypes.SearchResult;

public interface Backend {
    SearchResult query(String queryString);
}
