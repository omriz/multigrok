package com.github.omriz.multigrok.servers;

import com.github.omriz.multigrok.services.HelloService;
import org.wso2.msf4j.MicroservicesRunner;

public class Frontend {
    public static void main(String[] args) {
        new MicroservicesRunner()
                .deploy(new HelloService())
                .start();
    }
}
