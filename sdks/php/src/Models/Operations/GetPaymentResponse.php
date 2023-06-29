<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class GetPaymentResponse
{
	
    public string $contentType;
    
    /**
     * OK
     * 
     * @var ?\formance\stack\Models\Shared\PaymentResponse $paymentResponse
     */
	
    public ?\formance\stack\Models\Shared\PaymentResponse $paymentResponse = null;
    
	
    public int $statusCode;
    
	
    public ?\Psr\Http\Message\ResponseInterface $rawResponse = null;
    
	public function __construct()
	{
		$this->contentType = "";
		$this->paymentResponse = null;
		$this->statusCode = 0;
		$this->rawResponse = null;
	}
}
