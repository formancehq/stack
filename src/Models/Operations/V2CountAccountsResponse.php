<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class V2CountAccountsResponse
{
    /**
     * HTTP response content type for this operation
     * 
     * @var string $contentType
     */
	
    public string $contentType;
    
    /**
     * $headers
     * 
     * @var array<string, array<string>> $headers
     */
	
    public array $headers;
    
    /**
     * HTTP response status code for this operation
     * 
     * @var int $statusCode
     */
	
    public int $statusCode;
    
    /**
     * Raw HTTP response; suitable for custom response parsing
     * 
     * @var ?\Psr\Http\Message\ResponseInterface $rawResponse
     */
	
    public ?\Psr\Http\Message\ResponseInterface $rawResponse;
    
    /**
     * Error
     * 
     * @var ?\formance\stack\Models\Shared\V2ErrorResponse $v2ErrorResponse
     */
	
    public ?\formance\stack\Models\Shared\V2ErrorResponse $v2ErrorResponse = null;
    
	public function __construct()
	{
		$this->contentType = "";
		$this->headers = [];
		$this->statusCode = 0;
		$this->rawResponse = null;
		$this->v2ErrorResponse = null;
	}
}
