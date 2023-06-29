<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


/**
 * ListRunsResponse - List of workflow instances
 * 
 * @package formance\stack\Models\Shared
 * @access public
 */
class ListRunsResponse
{
    /**
     * $data
     * 
     * @var array<\formance\stack\Models\Shared\WorkflowInstance> $data
     */
	#[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('array<formance\stack\Models\Shared\WorkflowInstance>')]
    public array $data;
    
	public function __construct()
	{
		$this->data = [];
	}
}
