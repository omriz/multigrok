package com.github.omriz.multigrok.services;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;

@Path("/hello")
public class HelloService {

    @GET
    @Path("/{name}")
    public String hello(@PathParam("name") String name) {
        return "Hello " + name;
    }

}