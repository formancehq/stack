<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class GetInstanceResponse
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
     * The workflow instance
     * 
     * @var ?\formance\stack\Models\Shared\GetWorkflowInstanceResponse $getWorkflowInstanceResponse
     */
	
    public ?\formance\stack\Models\Shared\GetWorkflowInstanceResponse $getWorkflowInstanceResponse = null;
    
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
		$this->getWorkflowInstanceResponse = null;
		$this->statusCode = 0;
		$this->rawResponse = null;
	}
}