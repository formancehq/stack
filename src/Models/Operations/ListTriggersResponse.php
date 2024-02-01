<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class ListTriggersResponse
{
    /**
     * HTTP response content type for this operation
     * 
     * @var string $contentType
     */
	
    public string $contentType;
    
    /**
     * General error
     * 
     * @var ?\formance\stack\Models\Shared\Error $error
     */
	
    public ?\formance\stack\Models\Shared\Error $error = null;
    
    /**
     * List of triggers
     * 
     * @var ?\formance\stack\Models\Shared\ListTriggersResponse $listTriggersResponse
     */
	
    public ?\formance\stack\Models\Shared\ListTriggersResponse $listTriggersResponse = null;
    
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
    
	public function __construct()
	{
		$this->contentType = "";
		$this->error = null;
		$this->listTriggersResponse = null;
		$this->statusCode = 0;
		$this->rawResponse = null;
	}
}