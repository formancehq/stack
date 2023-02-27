package com.formancehq;

import static org.junit.Assert.assertTrue;

import com.formance.formance.*;
import com.formance.formance.api.ServerApi;
import com.formance.formance.ServerConfiguration;
import com.formance.formance.model.ConfigInfoResponse;
import org.junit.Test;
import retrofit2.Response;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;
import retrofit2.converter.scalars.ScalarsConverterFactory;

import java.io.IOException;
import java.util.HashMap;

/**
 * Unit test for simple App.
 */
public class AppTest 
{
    /**
     * Rigorous Test :-)
     */
    @Test
    public void shouldAnswerWithTrue() throws IOException {

        JSON json = new JSON();

        Retrofit.Builder adapterBuilder = new Retrofit
                .Builder()
                .baseUrl("http://localhost:3068")
                .addConverterFactory(ScalarsConverterFactory.create())
                .addConverterFactory(GsonConverterFactory.create(json.getGson()));

        ApiClient apiClient = new ApiClient();
        apiClient.setAdapterBuilder(adapterBuilder);
        ServerApi serverApi = apiClient.createService(ServerApi.class);

        System.out.println("TODO: Actualy, just checking SDK compilation and Client instanciation");
    }
}
