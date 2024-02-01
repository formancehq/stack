# DebitWalletResponse


## Fields

| Field                                                                                                        | Type                                                                                                         | Required                                                                                                     | Description                                                                                                  |
| ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------ |
| `contentType`                                                                                                | *string*                                                                                                     | :heavy_check_mark:                                                                                           | HTTP response content type for this operation                                                                |
| `debitWalletResponse`                                                                                        | [?\formance\stack\Models\Shared\DebitWalletResponse](../../Models/Shared/DebitWalletResponse.md)             | :heavy_minus_sign:                                                                                           | Wallet successfully debited as a pending hold                                                                |
| `statusCode`                                                                                                 | *int*                                                                                                        | :heavy_check_mark:                                                                                           | HTTP response status code for this operation                                                                 |
| `rawResponse`                                                                                                | [\Psr\Http\Message\ResponseInterface](https://www.php-fig.org/psr/psr-7/#33-psrhttpmessageresponseinterface) | :heavy_check_mark:                                                                                           | Raw HTTP response; suitable for custom response parsing                                                      |
| `walletsErrorResponse`                                                                                       | [?\formance\stack\Models\Shared\WalletsErrorResponse](../../Models/Shared/WalletsErrorResponse.md)           | :heavy_minus_sign:                                                                                           | Error                                                                                                        |